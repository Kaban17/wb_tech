const imagesContainer = document.getElementById('imagesContainer');
let trackedImages = JSON.parse(localStorage.getItem('trackedImages') || '[]');

function saveTrackedImages() {
    localStorage.setItem('trackedImages', JSON.stringify(trackedImages));
}

async function uploadImage() {
    const input = document.getElementById('imageInput');
    if (!input.files[0]) return alert('Select a file');

    const formData = new FormData();
    formData.append('image', input.files[0]);

    try {
        const response = await fetch('/upload', {
            method: 'POST',
            body: formData
        });
        const result = await response.json();
        trackedImages.push(result.id);
        saveTrackedImages();
        refreshImages();
        input.value = '';
    } catch (err) {
        console.error('Upload failed', err);
        alert('Upload failed');
    }
}

async function deleteImage(id) {
    if (!confirm('Delete this image?')) return;
    try {
        await fetch(`/image/${id}`, { method: 'DELETE' });
        trackedImages = trackedImages.filter(i => i !== id);
        saveTrackedImages();
        refreshImages();
    } catch (err) {
        console.error('Delete failed', err);
    }
}

async function refreshImages() {
    imagesContainer.innerHTML = '';
    for (const id of trackedImages) {
        try {
            const response = await fetch(`/image/${id}`);
            const img = await response.json();
            
            const card = document.createElement('div');
            card.className = 'image-card';
            
            let content = `
                <div class="status ${img.status}">${img.status.toUpperCase()}</div>
                <small>ID: ${img.id}</small>
            `;

            if (img.status === 'completed' && img.processed_path) {
                content = `<img src="/data/${img.processed_path}">` + content;
            } else {
                content = `<div style="width:200px; height:150px; background:#eee; display:flex; align-items:center; justify-content:center;">In processing...</div>` + content;
            }

            content += `<br><button onclick="deleteImage('${img.id}')">Delete</button>`;
            card.innerHTML = content;
            imagesContainer.appendChild(card);
        } catch (err) {
            console.error(`Failed to fetch image ${id}`, err);
        }
    }
}

// Poll for updates every 3 seconds
setInterval(refreshImages, 3000);
refreshImages();
