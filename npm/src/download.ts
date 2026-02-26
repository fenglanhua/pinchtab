import * as fs from 'fs';
import * as path from 'path';
import * as https from 'https';
import * as crypto from 'crypto';
import { detectPlatform, getBinaryName, getBinaryPath } from './platform';

const GITHUB_REPO = 'pinchtab/pinchtab';

// Read version from package.json at build time
function getVersion(): string {
  const pkgPath = path.join(__dirname, '..', 'package.json');
  const pkg = JSON.parse(fs.readFileSync(pkgPath, 'utf-8'));
  return pkg.version;
}

function fetchUrl(url: string, maxRedirects = 5): Promise<Buffer> {
  return new Promise((resolve, reject) => {
    const attemptFetch = (currentUrl: string, redirectsRemaining: number) => {
      let agent: https.Agent | undefined;

      // Proxy support for corporate environments
      const proxyUrl = process.env.HTTPS_PROXY || process.env.HTTP_PROXY;
      if (proxyUrl) {
        try {
          const proxy = new URL(proxyUrl);
          const proxyPort = proxy.port ? parseInt(proxy.port, 10) : 8080;
          agent = new https.Agent({
            host: proxy.hostname,
            port: proxyPort,
            keepAlive: true,
          });
        } catch (err) {
          console.warn(`Warning: Invalid proxy URL ${proxyUrl}, ignoring`);
        }
      }

      const request = https.get(currentUrl, agent ? { agent } : {}, (response) => {
        // Handle redirects (301, 302, 307, 308)
        if ([301, 302, 307, 308].includes(response.statusCode || 0)) {
          if (redirectsRemaining <= 0) {
            reject(new Error(`Too many redirects from ${currentUrl}`));
            return;
          }

          const redirectUrl = response.headers.location;
          if (!redirectUrl) {
            reject(new Error(`Redirect without location header from ${currentUrl}`));
            return;
          }

          // Consume the response stream to avoid memory leaks
          response.resume();
          attemptFetch(redirectUrl, redirectsRemaining - 1);
          return;
        }

        if (response.statusCode === 404) {
          reject(new Error(`Not found: ${currentUrl}`));
          return;
        }

        if (response.statusCode !== 200) {
          reject(new Error(`HTTP ${response.statusCode}: ${currentUrl}`));
          return;
        }

        const chunks: Buffer[] = [];
        response.on('data', (chunk) => chunks.push(chunk));
        response.on('end', () => resolve(Buffer.concat(chunks)));
        response.on('error', reject);
      });

      request.on('error', reject);
    };

    attemptFetch(url, maxRedirects);
  });
}

async function downloadChecksums(version: string): Promise<Map<string, string>> {
  const url = `https://github.com/${GITHUB_REPO}/releases/download/v${version}/checksums.txt`;

  try {
    const data = await fetchUrl(url);
    const checksums = new Map<string, string>();

    data
      .toString('utf-8')
      .trim()
      .split('\n')
      .forEach((line) => {
        const [hash, filename] = line.split(/\s+/);
        if (hash && filename) {
          checksums.set(filename.trim(), hash.trim());
        }
      });

    return checksums;
  } catch (err) {
    throw new Error(
      `Failed to download checksums: ${(err as Error).message}. ` +
        `Ensure v${version} is released on GitHub with checksums.txt`
    );
  }
}

function verifySHA256(filePath: string, expectedHash: string): boolean {
  const hash = crypto.createHash('sha256');
  const data = fs.readFileSync(filePath);
  hash.update(data);
  const actualHash = hash.digest('hex');
  return actualHash.toLowerCase() === expectedHash.toLowerCase();
}

async function downloadBinary(platform: any, version: string): Promise<void> {
  const binaryName = getBinaryName(platform);
  const binaryPath = getBinaryPath(binaryName, version);

  // Skip if already exists
  if (fs.existsSync(binaryPath)) {
    console.log(`✓ Pinchtab binary already present: ${binaryPath}`);
    return;
  }

  // Fetch checksums
  console.log(`Downloading Pinchtab ${version} for ${platform.os}-${platform.arch}...`);
  const checksums = await downloadChecksums(version);

  if (!checksums.has(binaryName)) {
    throw new Error(
      `Binary not found in checksums: ${binaryName}. ` +
        `Available: ${Array.from(checksums.keys()).join(', ')}`
    );
  }

  const expectedHash = checksums.get(binaryName)!;
  const downloadUrl = `https://github.com/${GITHUB_REPO}/releases/download/v${version}/${binaryName}`;

  // Ensure directory exists (unless using custom PINCHTAB_BINARY_PATH)
  const binDir = path.dirname(binaryPath);
  if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
  }

  // Download binary
  return new Promise((resolve, reject) => {
    console.log(`Downloading from ${downloadUrl}...`);

    const file = fs.createWriteStream(binaryPath);
    let redirectCount = 0;
    const maxRedirects = 5;

    const performDownload = (url: string) => {
      https
        .get(url, (response) => {
          // Handle redirects (301, 302, 307, 308)
          if ([301, 302, 307, 308].includes(response.statusCode || 0)) {
            if (redirectCount >= maxRedirects) {
              fs.unlink(binaryPath, () => {});
              reject(new Error(`Too many redirects downloading ${downloadUrl}`));
              return;
            }

            const redirectUrl = response.headers.location;
            if (!redirectUrl) {
              fs.unlink(binaryPath, () => {});
              reject(new Error(`Redirect without location header from ${url}`));
              return;
            }

            redirectCount++;
            response.resume(); // Consume response stream
            performDownload(redirectUrl);
            return;
          }

          if (response.statusCode !== 200) {
            fs.unlink(binaryPath, () => {});
            reject(new Error(`HTTP ${response.statusCode}: ${url}`));
            return;
          }

          response.pipe(file);

          file.on('finish', () => {
            file.close();

            // Verify checksum
            if (!verifySHA256(binaryPath, expectedHash)) {
              fs.unlink(binaryPath, () => {});
              reject(
                new Error(
                  `Checksum verification failed for ${binaryName}. ` +
                    `Download may be corrupted. Please try again.`
                )
              );
              return;
            }

            // Make executable
            fs.chmodSync(binaryPath, 0o755);
            console.log(`✓ Verified and installed: ${binaryPath}`);
            resolve();
          });

          file.on('error', (err) => {
            fs.unlink(binaryPath, () => {});
            reject(err);
          });
        })
        .on('error', reject);
    };

    performDownload(downloadUrl);
  });
}

export async function ensureBinary(): Promise<string> {
  const platform = detectPlatform();
  const version = getVersion();

  await downloadBinary(platform, version);

  const binaryName = getBinaryName(platform);
  return getBinaryPath(binaryName, version);
}
