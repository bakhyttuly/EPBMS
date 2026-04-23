document.getElementById("register-form").addEventListener("submit", async function (e) {
    e.preventDefault();

    const payload = {
        full_name: document.getElementById("full_name").value,
        email: document.getElementById("email").value,
        password: document.getElementById("password").value,
        role: document.getElementById("role").value
    };

    const response = await fetch("/auth/register", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(payload)
    });

    const data = await response.json();

    if (!response.ok) {
        alert(data.error || "Register failed");
        return;
    }

    alert("Registered successfully");
    window.location.href = "/";
});