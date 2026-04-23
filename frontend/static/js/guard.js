async function guardPage(allowedRoles) {
    const response = await fetch("/me", {
        credentials: "same-origin"
    });

    if (!response.ok) {
        window.location.href = "/";
        return;
    }

    const user = await response.json();

    if (!allowedRoles.includes(user.role)) {
        if (user.role === "performer") {
            window.location.href = "/my-schedule-page";
            return;
        }

        if (user.role === "admin" || user.role === "organizer") {
            window.location.href = "/dashboard-page";
            return;
        }

        window.location.href = "/";
    }
}