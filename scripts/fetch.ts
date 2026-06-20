import * as fs from 'fs';
import * as path from 'path';
import { AetherClient } from '../packages/sdk/src/client';

async function run() {
    const args = process.argv.slice(2);
    if (args.length < 1) {
        console.error("Usage: npx ts-node fetch.ts <capsule_uri_or_id> [output_file]");
        process.exit(1);
    }

    const uri = args[0];
    const outputFile = args[1] || 'downloaded_file.txt';

    const client = new AetherClient('http://localhost:5000');

    console.log(`Fetching ${uri} from the Aether network...`);
    try {
        const response = await client.fetchByUri(uri);
        if (!response.ok || !response.data) {
            console.error("Failed to fetch capsule or capsule is empty.");
            process.exit(1);
        }
        
        fs.writeFileSync(outputFile, response.data);
        console.log("=================================================");
        console.log(`SUCCESS! File downloaded to: ${outputFile}`);
        console.log("=================================================");
    } catch (e: any) {
        console.error("Failed to fetch:", e.message);
    }
}

run();
