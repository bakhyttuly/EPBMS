async function loadNavbar() {
    const response = await fetch("/me");
    if (!response.ok) return;

    const user = await response.json();

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
    }

    if (user.role === "organizer") {
        links = `
            <a href="/dashboard-page">Dashboard</a>
            <a href="/bookings-page">My Bookings</a>
            <a href="/calendar-page">Calendar</a>
        `;
    }

    if (user.role === "performer") {
        links = `
            <a href="/my-schedule-page">My Schedule</a>
        `;
    }

    nav.innerHTML = links + `<a href="/logout">Logout</a>`;
}

loadNavbar();