const axios = require('axios');

class AetherClient {
    constructor(config = {}) {
        this.daemonUrl = config.daemonUrl || 'http://localhost:5000';
    }

    /**
     * Check if the Aether Go daemon is online
     */
    async status() {
        try {
            const response = await axios.get(`${this.daemonUrl}/api/v1/status`);
            return response.data;
        } catch (error) {
            throw new Error(`Aether daemon is offline or unreachable at ${this.daemonUrl}. Please run the daemon first.`);
        }
    }

    /**
     * Publish an Envelope to the Aether network
     * @param {Object} envelope - The metadata envelope
     */
    async publish(envelope) {
        try {
            const response = await axios.post(`${this.daemonUrl}/api/v1/publish`, envelope);
            return response.data;
        } catch (error) {
            throw new Error(`Failed to publish envelope: ${error.message}`);
        }
    }

    /**
     * Upload a raw file to the daemon to be chunked into Boxo and get CIDs
     * @param {Buffer|Stream} fileData 
     */
    async uploadPayload(fileData) {
        try {
            const response = await axios.post(`${this.daemonUrl}/api/v1/payload`, fileData, {
                headers: { 'Content-Type': 'application/octet-stream' }
            });
            return response.data.cids;
        } catch (error) {
            throw new Error(`Failed to upload payload: ${error.message}`);
        }
    }
}

module.exports = { AetherClient };
