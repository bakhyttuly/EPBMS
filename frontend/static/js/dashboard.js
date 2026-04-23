async function loadDashboardStats() {
    try {
        const response = await authenticatedFetch("/api/v1/admin/stats");
        if (!response) return;

        const result = await response.json();
        if (!result.success) return;
        
        const data = result.data;

        document.getElementById("total-bookings").textContent = data.total_bookings ?? 0;
        document.getElementById("active-bookings").textContent = data.confirmed_bookings ?? 0;
        document.getElementById("completed-bookings").textContent = data.completed_bookings ?? 0;
        document.getElementById("performers-count").textContent = data.total_performers ?? 0;
    } catch (error) {
        console.error("Error loading stats:", error);
    }
}

loadDashboardStats();
setInterval(loadDashboardStats, 30000);
