document.addEventListener("DOMContentLoaded", () => {
  const searchInput = document.getElementById("search-input");
  const searchButton = document.getElementById("search-button");
  const showAllButton = document.getElementById("show-all-button");
  const commentsSection = document.getElementById("comments-section");
  const newCommentContent = document.getElementById("new-comment-content");
  const newCommentParentId = document.getElementById("new-comment-parent-id");
  const submitCommentButton = document.getElementById("submit-comment");

  const API_BASE_URL = "http://localhost:8001";

  // Function to fetch comments from the API
  async function fetchComments(params = {}) {
    let url = `${API_BASE_URL}/comments`; // Now targeting /comments
    const queryParams = new URLSearchParams();

    if (params.query) {
      queryParams.append("query", params.query);
    }
    if (params.parent) {
      queryParams.append("parent", params.parent);
    }
    // Add pagination/sorting params if needed

    if (queryParams.toString()) {
      url += `?${queryParams.toString()}`;
    }

    try {
      const response = await fetch(url);
      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(
          `HTTP error! status: ${response.status} - ${errorText}`,
        );
      }
      const data = await response.json();
      return data;
    } catch (error) {
      console.error("Error fetching comments:", error);
      commentsSection.innerHTML = `<p style="color: red;">Error loading comments: ${error.message}</p>`;
      return [];
    }
  }

  // Function to build a hierarchical tree from a flat list of comments
  function buildCommentTree(flatComments) {
    const commentsById = new Map();
    const rootComments = [];

    flatComments.forEach((comment) => {
      commentsById.set(comment.id, { ...comment, children: [] });
    });

    flatComments.forEach((comment) => {
      if (comment.parent_id !== null && commentsById.has(comment.parent_id)) {
        commentsById
          .get(comment.parent_id)
          .children.push(commentsById.get(comment.id));
      } else {
        rootComments.push(commentsById.get(comment.id));
      }
    });

    // Sort children by creation date
    const sortChildren = (comments) => {
      comments.sort((a, b) => new Date(a.created_at) - new Date(b.created_at));
      comments.forEach((comment) => {
        if (comment.children.length > 0) {
          sortChildren(comment.children);
        }
      });
    };
    sortChildren(rootComments);

    return rootComments;
  }

  // Function to render comments recursively
  function renderComment(comment, level) {
    const commentDiv = document.createElement("div");
    commentDiv.className = `comment-item comment-level-${level}`;
    commentDiv.dataset.id = comment.id;

    commentDiv.innerHTML = `
            <p class="comment-content">${comment.text}</p>
            <p class="comment-meta">ID: ${comment.id} - Author: ${comment.author_id} - Created: ${new Date(comment.created_at).toLocaleString()}</p>
            <div class="comment-actions">
                <button class="reply-button" data-id="${comment.id}">Reply</button>
                <button class="delete-button" data-id="${comment.id}">Delete</button>
            </div>
        `;
    commentsSection.appendChild(commentDiv);

    comment.children.forEach((child) => renderComment(child, level + 1));
  }

  // Main function to load and display comments
  async function loadAndDisplayComments(params = {}) {
    commentsSection.innerHTML = "<p>Loading comments...</p>";
    const flatComments = await fetchComments(params);
    commentsSection.innerHTML = ""; // Clear loading message

    if (flatComments.length === 0) {
      commentsSection.innerHTML = "<p>No comments found.</p>";
    } else {
      const commentTree = buildCommentTree(flatComments);
      commentTree.forEach((comment) => renderComment(comment, 0));
    }
  }

  // Event Listeners
  searchButton.addEventListener("click", () => {
    const searchTerm = searchInput.value.trim();
    loadAndDisplayComments({ query: searchTerm });
  });

  showAllButton.addEventListener("click", () => {
    searchInput.value = ""; // Clear search input
    loadAndDisplayComments(); // Fetch all comments (empty query)
  });

  searchInput.addEventListener("keypress", (event) => {
    if (event.key === "Enter") {
      const searchTerm = searchInput.value.trim();
      loadAndDisplayComments({ query: searchTerm });
    }
  });

  submitCommentButton.addEventListener("click", async () => {
    const content = newCommentContent.value.trim();
    const parentId = newCommentParentId.value.trim();

    if (!content) {
      alert("Comment content cannot be empty.");
      return;
    }

    const commentData = {
      author_id: 123, // Hardcoded for now
      text: content,
      parent_id: parentId ? parseInt(parentId, 10) : null,
    };

    try {
      const response = await fetch(`${API_BASE_URL}/comments`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(commentData),
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(
          `HTTP error! status: ${response.status} - ${errorText}`,
        );
      }

      newCommentContent.value = "";
      newCommentParentId.value = "";
      alert("Comment added successfully!");
      loadAndDisplayComments(); // Refresh comments
    } catch (error) {
      console.error("Error submitting comment:", error);
      alert(`Failed to add comment: ${error.message}`);
    }
  });

  commentsSection.addEventListener("click", async (event) => {
    if (event.target.classList.contains("delete-button")) {
      const commentId = event.target.dataset.id;
      if (
        confirm(
          `Are you sure you want to delete comment ID: ${commentId} and all its replies?`,
        )
      ) {
        try {
          const response = await fetch(
            `${API_BASE_URL}/comments/${commentId}`,
            {
              method: "DELETE",
            },
          );

          if (!response.ok) {
            const errorText = await response.text();
            throw new Error(
              `HTTP error! status: ${response.status} - ${errorText}`,
            );
          }

          alert("Comment deleted successfully!");
          loadAndDisplayComments(); // Refresh comments
        } catch (error) {
          console.error("Error deleting comment:", error);
          alert(`Failed to delete comment: ${error.message}`);
        }
      }
    } else if (event.target.classList.contains("reply-button")) {
      const commentId = event.target.dataset.id;
      newCommentParentId.value = commentId;
      newCommentContent.focus();
    }
  });

  // Initial load of comments
  loadAndDisplayComments();
});
