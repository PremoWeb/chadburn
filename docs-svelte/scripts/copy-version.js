#!/usr/bin/env node

import fs from 'fs';
import path from 'path';

// Get the project root directory (assuming this script is in docs-svelte/scripts)
const projectRoot = path.resolve(process.cwd(), '..');
const versionFilePath = path.join(projectRoot, 'VERSION');
const targetPath = path.join(process.cwd(), 'static', 'VERSION');

try {
  // Check if VERSION file exists
  if (fs.existsSync(versionFilePath)) {
    // Read the VERSION file
    const version = fs.readFileSync(versionFilePath, 'utf8').trim();
    console.log(`Found version: ${version}`);
    
    // Create a JSON file with the version
    const versionData = { version };
    fs.writeFileSync(
      path.join(process.cwd(), 'static', 'data', 'version.json'),
      JSON.stringify(versionData, null, 2)
    );
    console.log('Version data written to static/data/version.json');
    
    // Also copy the raw VERSION file
    fs.writeFileSync(targetPath, version);
    console.log('VERSION file copied to static/VERSION');
  } else {
    console.error('VERSION file not found in project root');
    process.exit(1);
  }
} catch (error) {
  console.error('Error copying VERSION file:', error);
  process.exit(1);
} 