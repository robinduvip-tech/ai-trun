const sharp = require('sharp')

const svg = `<svg xmlns="http://www.w3.org/2000/svg" width="1024" height="1024" viewBox="0 0 100 100" fill="none">
  <defs>
    <linearGradient id="bg" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#4F46E5"/>
      <stop offset="100%" stop-color="#7C3AED"/>
    </linearGradient>
  </defs>
  <!-- Rounded square background -->
  <rect x="6" y="6" width="88" height="88" rx="22" fill="url(#bg)"/>
  <!-- Connections -->
  <line x1="50" y1="50" x2="28" y2="26" stroke="white" stroke-width="3" stroke-linecap="round" opacity="0.45"/>
  <line x1="50" y1="50" x2="74" y2="30" stroke="white" stroke-width="3" stroke-linecap="round" opacity="0.40"/>
  <line x1="50" y1="50" x2="26" y2="68" stroke="white" stroke-width="3" stroke-linecap="round" opacity="0.40"/>
  <line x1="50" y1="50" x2="74" y2="72" stroke="white" stroke-width="3" stroke-linecap="round" opacity="0.45"/>
  <!-- Outer nodes -->
  <circle cx="28" cy="26" r="6" fill="white" opacity="0.85"/>
  <circle cx="74" cy="30" r="6" fill="white" opacity="0.75"/>
  <circle cx="26" cy="68" r="6" fill="white" opacity="0.75"/>
  <circle cx="74" cy="72" r="6" fill="white" opacity="0.85"/>
  <!-- Center node -->
  <circle cx="50" cy="50" r="9" fill="white" opacity="0.95"/>
  <!-- Orbit ring -->
  <circle cx="50" cy="50" r="20" stroke="white" stroke-width="2.5" stroke-dasharray="6 3 2 4" fill="none" opacity="0.35"/>
</svg>`

async function main() {
  await sharp(Buffer.from(svg))
    .resize(1024, 1024)
    .png()
    .toFile('F:\\workspace\\ai-trun\\desktop\\build\\appicon.png')
  console.log('appicon.png generated')
}

main().catch(e => { console.error(e); process.exit(1) })
