document.getElementById("login-form").addEventListener("submit", async function (e) {
    e.preventDefault();

    const payload = {
        email: document.getElementById("email").value,
        password: document.getElementById("password").value
    };

    const response = await fetch("/auth/login", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(payload)
    });

    const data = await response.json();

    if (!response.ok) {
        alert(data.error || "Login failed");
        return;
    }

    if (data.role === "performer") {
        window.location.href = "/my-schedule-page";
    } else {
        window.location.href = "/dashboard-page";
    }
});