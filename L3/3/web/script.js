document.addEventListener('DOMContentLoaded', () => {
    const commentsSection = document.getElementById('comments-section');
    const searchInput = document.getElementById('search-input');
    const searchButton = document.getElementById('search-button');
    const newCommentContent = document.getElementById('new-comment-content');
    const newCommentParentId = document.getElementById('new-comment-parent-id');
    const submitCommentButton = document.getElementById('submit-comment');

    const API_BASE_URL = 'http://localhost:8001'; // Assuming your Go API runs on 8001, as defined in .env

    async function fetchComments(parentId = '', searchTerm = '') {
        let url = `${API_BASE_URL}/comments`;
        const params = new URLSearchParams();

        if (parentId) {
            params.append('parent', parentId);
        }
        if (searchTerm) {
            params.append('query', searchTerm);
        }

        if (params.toString()) {
            url += `?${params.toString()}`;
        }

        try {
            const response = await fetch(url);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const comments = await response.json();
            displayComments(comments, 0);
        } catch (error) {
            console.error('Error fetching comments:', error);
            commentsSection.innerHTML = `<p>Error loading comments: ${error.message}</p>`;
        }
    }

    function displayComments(comments, level) {
        if (level === 0) {
            commentsSection.innerHTML = ''; // Clear existing comments only for top level
        }

        if (!comments || comments.length === 0) {
            if (level === 0) {
                commentsSection.innerHTML = '<p>No comments found.</p>';
            }
            return;
        }

        comments.forEach(comment => {
            const commentDiv = document.createElement('div');
            commentDiv.className = `comment-item level-${level}`;
            commentDiv.dataset.id = comment.id;

            let parentInfo = '';
            if (comment.parent_id) {
                parentInfo = `(Reply to: ${comment.parent_id})`;
            }

            commentDiv.innerHTML = `
                <p class="comment-content">${comment.content}</p>
                <p class="comment-meta">ID: ${comment.id} ${parentInfo} - Created: ${new Date(comment.created_at).toLocaleString()}</p>
                <div class="comment-actions">
                    <button class="delete-comment" data-id="${comment.id}">Delete</button>
                    <button class="reply-comment" data-id="${comment.id}">Reply</button>
                </div>
            `;
            commentsSection.appendChild(commentDiv);

            // Recursively display children
            if (comment.children && comment.children.length > 0) {
                displayComments(comment.children, level + 1);
            }
        });
    }

    searchButton.addEventListener('click', () => {
        const searchTerm = searchInput.value.trim();
        fetchComments('', searchTerm);
    });

    submitCommentButton.addEventListener('click', async () => {
        const content = newCommentContent.value.trim();
        const parentId = newCommentParentId.value.trim();

        if (!content) {
            alert('Comment content cannot be empty.');
            return;
        }

        try {
            const response = await fetch(`${API_BASE_URL}/comments`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    content: content,
                    parent_id: parentId ? parseInt(parentId) : null,
                }),
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(`HTTP error! status: ${response.status} - ${errorData.error}`);
            }

            newCommentContent.value = '';
            newCommentParentId.value = '';
            fetchComments(); // Refresh comments
        } catch (error) {
            console.error('Error submitting comment:', error);
            alert(`Failed to submit comment: ${error.message}`);
        }
    });

    commentsSection.addEventListener('click', async (event) => {
        if (event.target.classList.contains('delete-comment')) {
            const commentId = event.target.dataset.id;
            if (confirm(`Are you sure you want to delete comment ID: ${commentId} and all its replies?`)) {
                try {
                    const response = await fetch(`${API_BASE_URL}/comments/${commentId}`, {
                        method: 'DELETE',
                    });

                    if (!response.ok) {
                        const errorData = await response.json();
                        throw new Error(`HTTP error! status: ${response.status} - ${errorData.error}`);
                    }

                    fetchComments(); // Refresh comments
                } catch (error) {
                    console.error('Error deleting comment:', error);
                    alert(`Failed to delete comment: ${error.message}`);
                }
            }
        } else if (event.target.classList.contains('reply-comment')) {
            const commentId = event.target.dataset.id;
            newCommentParentId.value = commentId;
            newCommentContent.focus();
        }
    });

    // Initial load of comments
    fetchComments();
});
