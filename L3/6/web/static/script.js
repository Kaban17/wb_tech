document.addEventListener('DOMContentLoaded', () => {
    const API_BASE_URL = 'http://localhost:8001';

    // Item Form Elements
    const itemForm = document.getElementById('add-item-form');
    const itemIdInput = document.getElementById('itemId');
    const itemTypeInput = document.getElementById('item-type');
    const itemAmountInput = document.getElementById('item-amount');
    const itemDateInput = document.getElementById('item-date');
    const itemCategoryInput = document.getElementById('item-category');
    const itemDescriptionInput = document.getElementById('item-description');
    const saveItemButton = itemForm.querySelector('button[type="submit"]'); // Assuming the submit button is the save button
    const clearFormButton = document.getElementById('clear-form-button'); // Assuming there's a clear form button

    // Filter Form Elements
    const filterItemsForm = document.getElementById('filter-items-form');
    const filterTypeInput = document.getElementById('filter-type');
    const filterCategoryInput = document.getElementById('filter-category');
    const applyFilterButton = filterItemsForm.querySelector('button[type="submit"]');
    const clearFiltersButton = document.getElementById('clear-filters');

    // Items Table Elements
    const itemsTableBody = document.querySelector('#items-table tbody');

    // Analytics Form Elements
    const analyticsForm = document.getElementById('analytics-form');
    const analyticsFromDateInput = document.getElementById('analytics-from');
    const analyticsToDateInput = document.getElementById('analytics-to');
    const getAnalyticsButton = analyticsForm.querySelector('button[type="submit"]');
    const analyticsResultsDiv = document.getElementById('analytics-results');

    // --- Utility Functions ---
    function showMessage(message, isError = false) {
        const msgDiv = document.createElement('div');
        msgDiv.textContent = message;
        msgDiv.className = isError ? 'error-message' : 'success-message';
        // Append to a general message area, or just log for now
        console.log(message);
        if (isError) {
            console.error(message);
        }
        // You might want to display this in a dedicated message area on the page
    }

    function clearItemForm() {
        itemIdInput.value = '';
        itemTypeInput.value = 'income'; // Default value
        itemAmountInput.value = '';
        itemDateInput.value = '';
        itemCategoryInput.value = '';
        itemDescriptionInput.value = '';
        saveItemButton.textContent = 'Add Item';
    }

    // --- API Interaction Functions ---

    async function fetchAllItems(filters = {}) {
        let url = `${API_BASE_URL}/items`;
        const queryParams = new URLSearchParams();

        if (filters.type) {
            queryParams.append('type', filters.type);
        }
        if (filters.category) {
            queryParams.append('category', filters.category);
        }

        if (queryParams.toString()) {
            url += `?${queryParams.toString()}`;
        }

        try {
            const response = await fetch(url);
            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(`HTTP error! status: ${response.status} - ${errorText}`);
            }
            return await response.json();
        } catch (error) {
            showMessage(`Error fetching items: ${error.message}`, true);
            return [];
        }
    }

    async function fetchSingleItem(id) {
        try {
            const response = await fetch(`${API_BASE_URL}/items/${id}`);
            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(`HTTP error! status: ${response.status} - ${errorText}`);
            }
            return await response.json();
        } catch (error) {
            showMessage(`Error fetching item ${id}: ${error.message}`, true);
            return null;
        }
    }

    async function saveItem(itemData) {
        const method = itemData.id ? 'PUT' : 'POST';
        const url = itemData.id ? `${API_BASE_URL}/items/${itemData.id}` : `${API_BASE_URL}/items`;

        try {
            const response = await fetch(url, {
                method: method,
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(itemData),
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(`HTTP error! status: ${response.status} - ${errorText}`);
            }
            showMessage(`Item ${itemData.id ? 'updated' : 'added'} successfully!`);
            clearItemForm();
            loadItems();
        } catch (error) {
            showMessage(`Failed to save item: ${error.message}`, true);
        }
    }

    async function deleteItem(id) {
        if (!confirm(`Are you sure you want to delete item ID: ${id}?`)) {
            return;
        }
        try {
            const response = await fetch(`${API_BASE_URL}/items/${id}`, {
                method: 'DELETE',
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(`HTTP error! status: ${response.status} - ${errorText}`);
            }
            showMessage(`Item ${id} deleted successfully!`);
            loadItems();
        } catch (error) {
            showMessage(`Failed to delete item: ${error.message}`, true);
        }
    }

    async function getAnalytics(from, to) {
        const url = `${API_BASE_URL}/analytics?from=${encodeURIComponent(from)}&to=${encodeURIComponent(to)}`;
        try {
            const response = await fetch(url);
            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(`HTTP error! status: ${response.status} - ${errorText}`);
            }
            return await response.json();
        } catch (error) {
            showMessage(`Error fetching analytics: ${error.message}`, true);
            return null;
        }
    }

    // --- Rendering Functions ---

    function renderItemsTable(items) {
        itemsTableBody.innerHTML = '';
        if (items.length === 0) {
            itemsTableBody.innerHTML = '<tr><td colspan="7">No items found.</td></tr>';
            return;
        }

        items.forEach(item => {
            const row = itemsTableBody.insertRow();
            row.innerHTML = `
                <td>${item.id}</td>
                <td>${item.type}</td>
                <td>${item.amount.toFixed(2)}</td>
                <td>${new Date(item.date).toLocaleString()}</td>
                <td>${item.category}</td>
                <td>${item.description || ''}</td>
                <td class="action-buttons">
                    <button class="edit-button" data-id="${item.id}">Edit</button>
                    <button class="delete-button" data-id="${item.id}">Delete</button>
                </td>
            `;
        });
    }

    function renderAnalytics(analytics) {
        analyticsResultsDiv.innerHTML = '';
        if (!analytics) {
            analyticsResultsDiv.innerHTML = '<p>No analytics data available.</p>';
            return;
        }

        analyticsResultsDiv.innerHTML = `
            <p><strong>Total Sum:</strong> ${analytics.sum.toFixed(2)}</p>
            <p><strong>Average Amount:</strong> ${analytics.average.toFixed(2)}</p>
            <p><strong>Total Count:</strong> ${analytics.count}</p>
            <p><strong>Median:</strong> ${analytics.median.toFixed(2)}</p>
            <p><strong>90th Percentile:</strong> ${analytics.percentile90.toFixed(2)}</p>
        `;
    }

    // --- Main Load Function ---
    async function loadItems(filters = {}) {
        const items = await fetchAllItems(filters);
        renderItemsTable(items);
    }

    // --- Event Listeners ---

    itemForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const itemData = {
            id: itemIdInput.value ? parseInt(itemIdInput.value, 10) : 0,
            type: itemTypeInput.value,
            amount: parseFloat(itemAmountInput.value),
            date: new Date(itemDateInput.value).toISOString(), // Convert to ISO string for backend
            category: itemCategoryInput.value,
            description: itemDescriptionInput.value,
        };
        await saveItem(itemData);
    });

    clearFormButton.addEventListener('click', clearItemForm);

    filterItemsForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const filters = {
            type: filterTypeInput.value,
            category: filterCategoryInput.value,
        };
        await loadItems(filters);
    });

    clearFiltersButton.addEventListener('click', () => {
        filterTypeInput.value = '';
        filterCategoryInput.value = '';
        loadItems(); // Load all items
    });

    itemsTableBody.addEventListener('click', async (e) => {
        if (e.target.classList.contains('edit-button')) {
            const id = parseInt(e.target.dataset.id, 10);
            const item = await fetchSingleItem(id);
            if (item) {
                itemIdInput.value = item.id;
                itemTypeInput.value = item.type;
                itemAmountInput.value = item.amount;
                // Convert ISO string to YYYY-MM-DDTHH:MM for datetime-local input
                itemDateInput.value = new Date(item.date).toISOString().slice(0, 16);
                itemCategoryInput.value = item.category;
                itemDescriptionInput.value = item.description;
                saveItemButton.textContent = 'Update Item';
            }
        } else if (e.target.classList.contains('delete-button')) {
            const id = parseInt(e.target.dataset.id, 10);
            await deleteItem(id);
        }
    });

    analyticsForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const fromDate = analyticsFromDateInput.value;
        const toDate = analyticsToDateInput.value;

        if (!fromDate || !toDate) {
            showMessage('Please select both From and To dates for analytics.', true);
            return;
        }

        // Convert to ISO string for backend
        const fromISO = new Date(fromDate).toISOString();
        const toISO = new Date(toDate).toISOString();

        const analytics = await getAnalytics(fromISO, toISO);
        renderAnalytics(analytics);
    });

    // --- Initial Load ---
    loadItems();
    clearItemForm(); // Initialize form state
    // Set default dates for analytics (e.g., last 30 days)
    const today = new Date();
    const thirtyDaysAgo = new Date(today);
    thirtyDaysAgo.setDate(today.getDate() - 30);
    analyticsToDateInput.value = today.toISOString().slice(0, 10);
    analyticsFromDateInput.value = thirtyDaysAgo.toISOString().slice(0, 10);
});
