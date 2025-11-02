const userButton = document.getElementById("userButton");
const userMenu = document.getElementById("userMenu");
const userName = document.getElementById("userName");
const logoutBtn = document.getElementById("logoutBtn");

const loginForm = document.getElementById("loginForm");
const signupForm = document.getElementById("signupForm");
const switchAuth = document.getElementById("switchAuth");
const loginError = document.getElementById("loginError");
const signupError = document.getElementById("signupError");

const modal = new bootstrap.Modal(document.getElementById("authModal"));

const API_URL = "http://localhost:3003/authentification";

function logout() {
  localStorage.removeItem("jwt");
  localStorage.removeItem("user");
  location.reload();
}

// Handle login / signup toggle
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

// LOGIN
loginForm.addEventListener("submit", async (e) => {
  e.preventDefault();
  loginError.classList.add("d-none");
  const payload = {
    email: loginEmail.value.trim(),
    password: loginPassword.value.trim(),
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

// SIGNUP
signupForm.addEventListener("submit", async (e) => {
  e.preventDefault();
  signupError.classList.add("d-none");

  const payload = {
    role: "user",
    firstname: signupFirstname.value.trim(),
    lastname: signupLastname.value.trim(),
    email: signupEmail.value.trim(),
    phone: signupPhone.value.trim(),
    password: signupPassword.value.trim(),
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

// User button click handler
userButton.addEventListener("click", (e) => {
  e.preventDefault();
  const token = localStorage.getItem("jwt");
  if (!token) {
    modal.show();
  }
});

// Logout
logoutBtn.addEventListener("click", (e) => {
  e.preventDefault();
  localStorage.removeItem("jwt");
  localStorage.removeItem("user");
  location.reload();
});

// Update UI if user is logged in
function updateUserUI() {
  const user = JSON.parse(localStorage.getItem("user") || "null");
  if (user) {
    userName.textContent = `${user.firstname} ${user.lastname}`;
    userButton.querySelector("img").src = createInitialsAvatar(
      user.firstname,
      user.lastname
    );
  } else {
    userName.textContent = "Гость";
    userButton.querySelector("img").src =
      "https://img.freepik.com/premium-vector/default-avatar-profile-icon-social-media-user-image-gray-avatar-icon-blank-profile-silhouette-vector-illustration_561158-3485.jpg?semt=ais_hybrid&w=740&q=80";
  }
}

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
