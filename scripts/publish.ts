import * as fs from 'fs';
import * as path from 'path';
import { AetherClient } from '../packages/sdk/src/client';

async function run() {
    const args = process.argv.slice(2);
    if (args.length < 2) {
        console.error("Usage: npx ts-node publish.ts <username> <file_path>");
        process.exit(1);
    }

    const username = args[0];
    const filePath = args[1];

    if (!fs.existsSync(filePath)) {
        console.error(`File not found: ${filePath}`);
        process.exit(1);
    }

    const payload = fs.readFileSync(filePath, 'utf-8');
    const client = new AetherClient('http://localhost:5000');

    console.log(`Publishing ${filePath} as ${username}...`);
    try {
        const response = await client.publish(username, payload, { persistent: true, visibility: 'public' });
        console.log("=================================================");
        console.log(`SUCCESS! Published Capsule ID:`);
        console.log(`   aether://${response.capsule_id}`);
        console.log("=================================================");
        console.log(`Tell your friend to run:`);
        console.log(`   npx ts-node fetch.ts aether://${response.capsule_id}`);
    } catch (e: any) {
        console.error("Failed to publish:", e.message);
    }
}

run();
