import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router";

type AuthMode = "sign-up" | "sign-in";

const AuthForm: React.FC = () => {
  const [formData, setFormData] = useState({
    email: "",
    password: "",
    firstName: "",
    lastName: "",
    address: "",
    phoneNumber: "",
  });
  const [authMode, setAuthMode] = useState<AuthMode>("sign-in"); // Стартовый режим: вход
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null); // Состояние для успешных сообщений
  const navigate = useNavigate();

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (token) {
      navigate("/profile"); // Если токен есть, перенаправляем в профиль
    }
  }, [navigate]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const isSignUpMode = authMode === "sign-up";

  const validatePassword = (password: string): string | null => {
    if (password.length < 8) {
      return "Password must be at least 8 characters long.";
    }

    const hasUpper = /[A-Z]/.test(password);
    const hasLower = /[a-z]/.test(password);
    const hasDigit = /[0-9]/.test(password);
    const hasSpecial = /[#!@$%^&*-]/.test(password);

    if (!hasUpper) {
      return "Password must contain at least one uppercase letter.";
    }
    if (!hasLower) {
      return "Password must contain at least one lowercase letter.";
    }
    if (!hasDigit) {
      return "Password must contain at least one digit.";
    }
    if (!hasSpecial) {
      return "Password must contain at least one special character.";
    }

    return null;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (isSignUpMode) {
      const passwordError = validatePassword(formData.password);
      if (passwordError) {
        setErrorMessage(passwordError);
        setTimeout(() => setErrorMessage(null), 5000);
        return;
      }
    }

    const endpoint = isSignUpMode ? "/auth/sign-up" : "/auth/sign-in";
    const method = "POST";
    const requestBody = {
      email: formData.email,
      password: formData.password,
      ...(isSignUpMode && {
        firstName: formData.firstName,
        lastName: formData.lastName,
        address: formData.address,
        phoneNumber: formData.phoneNumber,
      }),
    };

    try {
      const response = await fetch(`http://localhost:8088${endpoint}`, {
        method,
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(requestBody),
      });

      if (response.ok) {
        const data = await response.json();

        if (isSignUpMode && response.status === 201) {
          setSuccessMessage("Account successfully created");
          setTimeout(() => setSuccessMessage(null), 5000);
        }

        if (response.status !== 201 && data) {
          localStorage.setItem("token", data); // Сохраняем токен в localStorage
          navigate("/profile"); // Перенаправляем в профиль
        } else {
          setErrorMessage("Token not received");
          setTimeout(() => setErrorMessage(null), 5000);
        }
      } else {
        const errorData = await response.json();
        if (errorData.message) {
          console.log(errorData);
          setErrorMessage(errorData.message);
        } else {
          setErrorMessage("An error occurred");
        }
        setTimeout(() => setErrorMessage(null), 5000);
      }
    } catch (error) {
      setErrorMessage("An error occurred");
      setTimeout(() => setErrorMessage(null), 5000);
    }
  };

  return (
      <div className="auth-container">
        <h1>{isSignUpMode ? "Register" : "Sign In"}</h1>
        <form onSubmit={handleSubmit} className="auth-form">
          <label>Email</label>
          <input
              type="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              required
          />

          <label>Password</label>
          <input
              type="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              required
          />

          {isSignUpMode && (
              <>
                <label>First Name</label>
                <input
                    type="text"
                    name="firstName"
                    value={formData.firstName}
                    onChange={handleChange}
                />

                <label>Last Name</label>
                <input
                    type="text"
                    name="lastName"
                    value={formData.lastName}
                    onChange={handleChange}
                />

                <label>Address</label>
                <input
                    type="text"
                    name="address"
                    value={formData.address}
                    onChange={handleChange}
                />

                <label>Phone Number</label>
                <input
                    type="text"
                    name="phoneNumber"
                    value={formData.phoneNumber}
                    onChange={handleChange}
                />
              </>
          )}

          <button type="submit">{isSignUpMode ? "Register" : "Sign In"}</button>
        </form>

        {errorMessage && (
            <div className="error-message">
              <p>{errorMessage}</p> {/* Выводим сообщение об ошибке */}
            </div>
        )}
        {successMessage && (
            <div className="success-message">
              <p>{successMessage}</p>
            </div>
        )}

        <button
            onClick={() => setAuthMode(isSignUpMode ? "sign-in" : "sign-up")}
            className="switch-auth-mode"
        >
          {isSignUpMode
              ? "Already have an account? Sign In"
              : "Don't have an account? Register"}
        </button>
      </div>
  );
};

export default AuthForm;