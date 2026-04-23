function loadNavbar() {
    const userStr = localStorage.getItem("user");
    if (!userStr) return;

    const user = JSON.parse(userStr);
    const nav = document.querySelector(".sidebar nav");
    if (!nav) return;

    let links = "";

    if (user.role === "admin") {
        links = `
            <a href="/dashboard-page">Dashboard</a>
            <a href="/performers-page">Performers</a>
            <a href="/bookings-page">Bookings</a>
            <a href="/calendar-page">Calendar</a>
        `;
    } else if (user.role === "client") {
        links = `
            <a href="/performers-page">Browse Performers</a>
            <a href="/bookings-page">My Bookings</a>
        `;
    } else if (user.role === "performer") {
        links = `
            <a href="/my-schedule-page">My Schedule</a>
        `;
    }

    nav.innerHTML = links + `<a href="#" id="logout-link">Logout</a>`;
    
    const logoutLink = document.getElementById("logout-link");
    if (logoutLink) {
        logoutLink.addEventListener("click", (e) => {
            e.preventDefault();
            localStorage.removeItem("token");
            localStorage.removeItem("user");
            window.location.href = "/";
        });
    }
}

loadNavbar();
