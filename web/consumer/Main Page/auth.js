// auth.js
document.addEventListener('DOMContentLoaded', function() {
    // Initialize elements
    const userButton = document.getElementById("userButton");
    const userName = document.getElementById("userName");
    const logoutBtn = document.getElementById("logoutBtn");

    const loginForm = document.getElementById("loginForm");
    const signupForm = document.getElementById("signupForm");
    const switchAuth = document.getElementById("switchAuth");
    const loginError = document.getElementById("loginError");
    const signupError = document.getElementById("signupError");

    const modalElement = document.getElementById("authModal");
    
    // Check if all elements exist
    if (!modalElement) {
        console.error('Auth modal element not found!');
        return;
    }

    const modal = new bootstrap.Modal(modalElement);
    const API_URL = "http://localhost:3004/authentification";

    // Handle login / signup toggle
    if (switchAuth) {
        switchAuth.addEventListener("click", (e) => {
            e.preventDefault();
            const isLoginVisible = !loginForm.classList.contains("d-none");
            loginForm.classList.toggle("d-none", isLoginVisible);
            signupForm.classList.toggle("d-none", !isLoginVisible);
            switchAuth.textContent = isLoginVisible
                ? "Уже есть аккаунт? Войти"
                : "Нет аккаунта? Зарегистрируйтесь";
            document.getElementById("authModalLabel").textContent = isLoginVisible
                ? "Регистрация"
                : "Вход в аккаунт";
        });
    }

    // LOGIN
    if (loginForm) {
        loginForm.addEventListener("submit", async (e) => {
            e.preventDefault();
            loginError.classList.add("d-none");
            const payload = {
                email: document.getElementById('loginEmail').value.trim(),
                password: document.getElementById('loginPassword').value.trim(),
            };

            try {
                const res = await fetch(`${API_URL}/login`, {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(payload),
                });
                const data = await res.json();

                if (!res.ok) throw new Error(data.error || "Ошибка входа");

                // Save token & user info
                localStorage.setItem("jwt", data.token);
                localStorage.setItem("user", JSON.stringify(data.user));

                modal.hide();
                updateUserUI();
            } catch (err) {
                loginError.textContent = err.message;
                loginError.classList.remove("d-none");
            }
        });
    }

    // SIGNUP
    if (signupForm) {
        signupForm.addEventListener("submit", async (e) => {
            e.preventDefault();
            signupError.classList.add("d-none");

            const payload = {
                role: "user",
                firstname: document.getElementById('signupFirstname').value.trim(),
                lastname: document.getElementById('signupLastname').value.trim(),
                email: document.getElementById('signupEmail').value.trim(),
                phone: document.getElementById('signupPhone').value.trim(),
                password: document.getElementById('signupPassword').value.trim(),
            };

            try {
                const res = await fetch(`${API_URL}/signup`, {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(payload),
                });
                const data = await res.json();

                if (!res.ok) throw new Error(data.error || "Ошибка регистрации");

                alert("Аккаунт создан! Теперь войдите.");
                signupForm.reset();
                switchAuth.click(); // Switch back to login
            } catch (err) {
                signupError.textContent = err.message;
                signupError.classList.remove("d-none");
            }
        });
    }

    // User button click handler - FIXED: prevent dropdown when showing modal
    if (userButton) {
        userButton.addEventListener("click", (e) => {
            const token = localStorage.getItem("jwt");
            if (!token) {
                e.preventDefault();
                e.stopPropagation();
                modal.show();
            }
            // If token exists, Bootstrap dropdown will handle the menu
        });
    }

    // Logout
    if (logoutBtn) {
        logoutBtn.addEventListener("click", (e) => {
            e.preventDefault();
            localStorage.removeItem("jwt");
            localStorage.removeItem("user");
            location.reload();
        });
    }

    // Update UI if user is logged in
    function updateUserUI() {
        const user = JSON.parse(localStorage.getItem("user") || "null");
        const userAvatar = document.getElementById('userAvatar');
        
        if (user) {
            if (userName) {
                userName.textContent = `${user.firstname} ${user.lastname}`;
            }
            if (userAvatar) {
                userAvatar.src = createInitialsAvatar(
                    user.firstname,
                    user.lastname
                );
            }
        } else {
            if (userName) {
                userName.textContent = "Гость";
            }
            if (userAvatar) {
                userAvatar.src = "https://img.freepik.com/premium-vector/default-avatar-profile-icon-social-media-user-image-gray-avatar-icon-blank-profile-silhouette-vector-illustration_561158-3485.jpg?semt=ais_hybrid&w=740&q=80";
            }
        }
    }

    // Initialize user UI on page load
    updateUserUI();

    function createInitialsAvatar(
        firstName,
        lastName,
        size = 32,
        bgColor = "#007bff",
        textColor = "#fff"
    ) {
        const canvas = document.createElement("canvas");
        canvas.width = size;
        canvas.height = size;
        const ctx = canvas.getContext("2d");

        // Draw background circle
        ctx.fillStyle = bgColor;
        ctx.beginPath();
        ctx.arc(size / 2, size / 2, size / 2, 0, Math.PI * 2);
        ctx.fill();

        // Draw initials
        ctx.fillStyle = textColor;
        ctx.font = `${size / 2}px sans-serif`;
        ctx.textAlign = "center";
        ctx.textBaseline = "middle";

        const initials = `${firstName?.[0] || ""}${
            lastName?.[0] || ""
        }`.toUpperCase();
        ctx.fillText(initials, size / 2, size / 2);

        return canvas.toDataURL();
    }
});