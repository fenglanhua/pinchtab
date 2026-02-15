// Stealth script injected into every new document via page.AddScriptToEvaluateOnNewDocument.
// Hides automation indicators from bot detection systems.

// Hide webdriver flag
Object.defineProperty(navigator, 'webdriver', { get: () => undefined });

// Fake chrome.runtime (automation detection)
if (!window.chrome) { window.chrome = {}; }
if (!window.chrome.runtime) { window.chrome.runtime = {}; }

// Fix notification permissions query
const originalQuery = window.navigator.permissions.query;
window.navigator.permissions.query = (parameters) => (
	parameters.name === 'notifications' ?
		Promise.resolve({ state: Notification.permission }) :
		originalQuery(parameters)
);

// Hide plugins length (headless has 0)
Object.defineProperty(navigator, 'plugins', {
	get: () => [1, 2, 3, 4, 5],
});

// Standard languages (headless has empty)
Object.defineProperty(navigator, 'languages', {
	get: () => ['en-GB', 'en-US', 'en'],
});
