// Emoji Colorization System
// Wraps emoji characters with spans that apply CSS filter colors

(function() {
  'use strict';

  // Define emoji patterns and their default color classes
  const emojiConfig = {
    '🦀': 'emoji-primary',    // Crab - primary brand color
    '📁': 'emoji-orange',     // Folder
    '🖥️': 'emoji-blue',       // Monitor
    '⚠️': 'emoji-warning',    // Warning
    '✓': 'emoji-success',     // Check
    '✗': 'emoji-error',       // X
    '🔒': 'emoji-muted',      // Lock
    '🔓': 'emoji-success',    // Unlock
    '📊': 'emoji-blue',       // Chart
    '🎯': 'emoji-primary',    // Target
    '⏱️': 'emoji-purple',     // Timer
    '💾': 'emoji-blue',       // Save
    '🔄': 'emoji-orange',     // Refresh
    '❌': 'emoji-error',      // Close
  };

  // Regular expression to match all configured emojis
  const emojiRegex = new RegExp(
    '(' + Object.keys(emojiConfig).map(e => e.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')).join('|') + ')',
    'g'
  );

  // Function to wrap an emoji with a colored span
  function wrapEmoji(emoji, className = '') {
    const defaultClass = emojiConfig[emoji] || 'emoji-muted';
    return `<span class="emoji ${className || defaultClass}">${emoji}</span>`;
  }

  // Function to colorize all emojis in a text string
  function colorizeEmojisInText(text, defaultClassName = '') {
    return text.replace(emojiRegex, (match) => {
      return wrapEmoji(match, defaultClassName);
    });
  }

  // Function to colorize all emojis in an element
  function colorizeEmojisInElement(element, defaultClassName = '') {
    // Process text nodes
    const walker = document.createTreeWalker(
      element,
      NodeFilter.SHOW_TEXT,
      {
        acceptNode: function(node) {
          // Skip if parent is already an emoji span
          if (node.parentElement && node.parentElement.classList.contains('emoji')) {
            return NodeFilter.FILTER_REJECT;
          }
          // Check if contains emoji
          return emojiRegex.test(node.textContent) ? 
            NodeFilter.FILTER_ACCEPT : NodeFilter.FILTER_REJECT;
        }
      }
    );

    const textNodes = [];
    while (walker.nextNode()) {
      textNodes.push(walker.currentNode);
    }

    // Process each text node
    textNodes.forEach(node => {
      const span = document.createElement('span');
      span.innerHTML = colorizeEmojisInText(node.textContent, defaultClassName);
      node.parentNode.replaceChild(span, node);
    });
  }

  // Function to colorize all emojis in the entire document
  function colorizeAllEmojis(defaultClassName = '') {
    colorizeEmojisInElement(document.body, defaultClassName);
  }

  // Helper to add color class to existing emoji elements
  function addColorToEmoji(element, colorClass) {
    if (!element.classList.contains('emoji')) {
      element.classList.add('emoji');
    }
    // Remove any existing color classes
    Object.values(emojiConfig).forEach(cls => {
      element.classList.remove(cls);
    });
    element.classList.remove('emoji-primary', 'emoji-blue', 'emoji-red', 
      'emoji-green', 'emoji-purple', 'emoji-orange', 'emoji-yellow', 
      'emoji-success', 'emoji-warning', 'emoji-error', 'emoji-muted');
    // Add new color class
    element.classList.add(colorClass);
  }

  // Mutation observer to handle dynamically added content
  function observeNewEmojis() {
    const observer = new MutationObserver((mutations) => {
      mutations.forEach((mutation) => {
        mutation.addedNodes.forEach((node) => {
          if (node.nodeType === Node.ELEMENT_NODE) {
            colorizeEmojisInElement(node);
          }
        });
      });
    });

    observer.observe(document.body, {
      childList: true,
      subtree: true
    });
  }

  // Auto-initialize on DOMContentLoaded
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', () => {
      colorizeAllEmojis();
      observeNewEmojis();
    });
  } else {
    // DOM already loaded
    colorizeAllEmojis();
    observeNewEmojis();
  }

  // Export utilities for manual use
  window.emojiColors = {
    wrap: wrapEmoji,
    colorizeText: colorizeEmojisInText,
    colorizeElement: colorizeEmojisInElement,
    colorizeAll: colorizeAllEmojis,
    addColorToEmoji: addColorToEmoji,
    config: emojiConfig
  };

})();