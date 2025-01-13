import fs from 'fs-extra';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

function copyAssets() {
  // First, copy assets to dist/assets
  const srcDir = path.join(__dirname, '../src/assets');
  const distDir = path.join(__dirname, '../dist/assets');
  const iosDir = path.join(__dirname, '../ios/App/App/public/assets');

  // Ensure directories exist
  fs.ensureDirSync(distDir);
  fs.ensureDirSync(iosDir);

  // Files to copy
  const filesToCopy = ['game.wasm', 'wasm_exec.js'];
  
  // Copy to dist
  filesToCopy.forEach(file => {
    fs.copySync(
      path.join(srcDir, file),
      path.join(distDir, file),
      { overwrite: true }
    );
    console.log(`Copied ${file} to dist/assets`);
  });

  // Copy to iOS public
  filesToCopy.forEach(file => {
    fs.copySync(
      path.join(srcDir, file),
      path.join(iosDir, file),
      { overwrite: true }
    );
    console.log(`Copied ${file} to iOS public assets`);
  });
}

try {
  copyAssets();
  console.log('Assets copied successfully');
} catch (err) {
  console.error('Error copying assets:', err);
  process.exit(1);
}
