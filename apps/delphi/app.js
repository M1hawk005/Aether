document.addEventListener('DOMContentLoaded', () => {
    let currentUser = "peer-a";

    // Try to get current peer id from status
    fetch('/api/status')
        .then(res => res.json())
        .then(data => {
            if (data.peer_id) currentUser = data.peer_id;
            loadFeed();
        })
        .catch(err => {
            console.error(err);
            loadFeed();
        });

    document.getElementById('btn-publish').addEventListener('click', () => {
        document.getElementById('publish-modal').classList.remove('hidden');
    });

    document.getElementById('btn-cancel').addEventListener('click', () => {
        document.getElementById('publish-modal').classList.add('hidden');
    });

    document.getElementById('publish-form').addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const title = document.getElementById('pub-title').value;
        const magnet = document.getElementById('pub-magnet').value;
        const desc = document.getElementById('pub-desc').value;

        const payload = JSON.stringify({
            app: 'tracker',
            title,
            magnet,
            description: desc,
            timestamp: Date.now()
        });

        // Publish to Aether
        try {
            const res = await fetch('/api/publish', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    publisher: currentUser,
                    payload: payload,
                    visibility: 'public',
                    persistent: true
                })
            });

            if (res.ok) {
                alert("Torrent Published to Aether!");
                document.getElementById('publish-modal').classList.add('hidden');
                document.getElementById('publish-form').reset();
                setTimeout(loadFeed, 500); // give daemon time to cache
            } else {
                const err = await res.json();
                alert("Failed to publish: " + err.error);
            }
        } catch (err) {
            console.error(err);
            alert("Error connecting to daemon");
        }
    });

    // Search functionality
    document.getElementById('search-input').addEventListener('input', (e) => {
        const query = e.target.value.toLowerCase();
        const cards = document.querySelectorAll('.torrent-card');
        cards.forEach(card => {
            const title = card.querySelector('.torrent-title').textContent.toLowerCase();
            const desc = card.querySelector('.torrent-desc').textContent.toLowerCase();
            if (title.includes(query) || desc.includes(query)) {
                card.style.display = 'block';
            } else {
                card.style.display = 'none';
            }
        });
    });
});

async function loadFeed() {
    const list = document.getElementById('torrent-list');
    list.innerHTML = '<div class="loading">Loading decentralized feed...</div>';

    try {
        const res = await fetch('/api/feed');
        const data = await res.json();
        
        if (!data.capsules || data.capsules.length === 0) {
            list.innerHTML = '<div class="loading">No torrents found on the network.</div>';
            return;
        }

        list.innerHTML = '';
        
        // Fetch payloads
        // Sort capsules by time
        data.capsules.sort((a,b) => new Date(b.created_at) - new Date(a.created_at));

        for (const cap of data.capsules) {
            try {
                const payloadRes = await fetch('/api/fetch/' + cap.capsule_id);
                if (payloadRes.ok) {
                    const payloadData = await payloadRes.json();
                    let parsed;
                    try {
                        parsed = JSON.parse(payloadData.data);
                    } catch (e) { continue; }
                    
                    if (parsed.app === 'tracker') {
                        renderTorrent(parsed, cap.capsule_id, cap.created_at);
                    }
                }
            } catch (err) {
                console.error("Failed to fetch capsule", cap.capsule_id, err);
            }
        }
        
        if (list.children.length === 0) {
            list.innerHTML = '<div class="loading">No tracker capsules found yet.</div>';
        }
    } catch (err) {
        list.innerHTML = '<div class="loading" style="color: #ef4444;">Failed to load feed. Is Aether daemon running?</div>';
    }
}

function renderTorrent(torrent, id, createdAt) {
    const list = document.getElementById('torrent-list');
    
    const card = document.createElement('div');
    card.className = 'torrent-card glass-panel';
    
    // Convert tracker URL
    const aetherTrackerUrl = `http://127.0.0.1:5000/announce`;
    
    // Safely append tracker to magnet link
    let finalMagnet = escapeHTML(torrent.magnet);
    if (!finalMagnet.includes('tr=')) {
        finalMagnet += `&tr=${encodeURIComponent(aetherTrackerUrl)}`;
    }
    
    card.innerHTML = `
        <div class="torrent-header">
            <div class="torrent-title">${escapeHTML(torrent.title)}</div>
        </div>
        <div class="torrent-meta">
            <span>Published: ${new Date(createdAt).toLocaleString()}</span>
            <span>Aether Capsule: ${id.substring(0, 8)}...</span>
        </div>
        <div class="torrent-desc">${escapeHTML(torrent.description)}</div>
        <a href="${finalMagnet}" class="magnet-link">Download (Magnet)</a>
    `;
    
    list.appendChild(card);
}

function escapeHTML(str) {
    if (!str) return '';
    return str.replace(/[&<>'"]/g, 
        tag => ({
            '&': '&amp;',
            '<': '&lt;',
            '>': '&gt;',
            "'": '&#39;',
            '"': '&quot;'
        }[tag] || tag)
    );
}
