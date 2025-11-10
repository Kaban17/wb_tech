document.addEventListener("DOMContentLoaded", function () {
  const searchInput = document.getElementById("search-input");
  const searchButton = document.getElementById("search-button");
  const commentsSection = document.getElementById("comments-section");
  const newCommentContent = document.getElementById("new-comment-content");
  const newCommentParentId = document.getElementById("new-comment-parent-id");
  const submitCommentButton = document.getElementById("submit-comment");

  // Function to fetch and display comments
  async function fetchAndDisplayComments(query = "") {
    commentsSection.innerHTML = ""; // Clear previous results
    let url = `/comments/search?query=${encodeURIComponent(query)}`;

    try {
      const response = await fetch(url);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const comments = await response.json();

      if (comments && comments.length > 0) {
        comments.forEach((comment) => {
          const commentDiv = document.createElement("div");
          commentDiv.classList.add("comment-item");
          commentDiv.innerHTML = `
                        <p><strong>ID:</strong> ${comment.id}</p>
                        <p>${comment.text}</p>
                        <p class="meta">Author: ${comment.author_id} | Created: ${new Date(comment.created_at).toLocaleString()}</p>
                        ${comment.parent_id ? `<p class="meta">Parent ID: ${comment.parent_id}</p>` : ""}
                    `;
          commentsSection.appendChild(commentDiv);
        });
      } else {
        commentsSection.innerHTML = "<p>No comments found.</p>";
      }
    } catch (error) {
      console.error("Error fetching comments:", error);
      commentsSection.innerHTML = `<p style="color: red;">Error fetching comments: ${error.message}</p>`;
    }
  }

  // Event listener for search button and Enter key
  searchButton.addEventListener("click", () =>
    fetchAndDisplayComments(searchInput.value),
  );
  searchInput.addEventListener("keypress", function (event) {
    if (event.key === "Enter") {
      fetchAndDisplayComments(searchInput.value);
    }
  });

  // Event listener for submitting a new comment
  submitCommentButton.addEventListener("click", async () => {
    const text = newCommentContent.value;
    const parentId = newCommentParentId.value;

    if (!text) {
      alert("Comment content cannot be empty.");
      return;
    }

    const commentData = {
      author_id: 123, // Hardcoded for now, ideally from user session
      text: text,
    };

    if (parentId) {
      commentData.parent_id = parseInt(parentId, 10);
    }

    try {
      const response = await fetch("/comments", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(commentData),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      newCommentContent.value = "";
      newCommentParentId.value = "";
      fetchAndDisplayComments(); // Refresh comments after adding a new one
      alert("Comment added successfully!");
    } catch (error) {
      console.error("Error adding comment:", error);
      alert(`Error adding comment: ${error.message}`);
    }
  });

  // Initial load of comments (fetch all comments with an empty query)
  fetchAndDisplayComments();
});
