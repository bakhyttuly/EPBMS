async function loadDashboardStats() {
    const response = await fetch("/dashboard/stats");
    const data = await response.json();

    document.getElementById("total-bookings").textContent = data.total_bookings ?? 0;
    document.getElementById("active-bookings").textContent = data.active_bookings ?? 0;
    document.getElementById("completed-bookings").textContent = data.completed_bookings ?? 0;
    document.getElementById("performers-count").textContent = data.performers_count ?? 0;
}

loadDashboardStats();