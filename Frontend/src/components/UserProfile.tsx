import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router";

interface User {
  first_name: string;
  last_name: string;
  address: string;
  phone_number: string;
  email: string;
}

interface Order {
  id: number;
  type: string;
  brand: string;
  model: string;
  year: number;
  image: string;
  total_price: number;
  start_date: string;
  end_date: string;
}

const UserProfile: React.FC = () => {
  const [user, setUser] = useState<User | null>(null);
  const [orders, setOrders] = useState<Order[]>([]);
  const [isEditing, setIsEditing] = useState(false);
  const [formData, setFormData] = useState<User>({
    first_name: "",
    last_name: "",
    address: "",
    phone_number: "",
    email: "",
  });

  const navigate = useNavigate();

  useEffect(() => {
    const token = localStorage.getItem("token");

    // Проверка токена
    if (!token) {
      navigate("/auth");
    } else {
      fetchUserProfile(token);
      fetchOrderHistory(token);
    }
  }, [navigate]);

  const fetchUserProfile = async (token: string) => {
    try {
      const response = await fetch("http://localhost:8088/profile/user", {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (response.status == 401) {
        localStorage.clear();
      }

      if (!response.ok) {
        throw new Error("Failed to fetch user data");
      }

      const data = await response.json();
      setUser(data);
      setFormData(data);
    } catch (error) {
      console.error(error);
    }
  };

  const fetchOrderHistory = async (token: string) => {
    try {
      const response = await fetch("http://localhost:8088/products/info", {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error("Failed to fetch order history");
      }

      const data = await response.json();
      setOrders(data);
    } catch (error) {
      console.error(error);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const token = localStorage.getItem("token");

    try {
      const response = await fetch("http://localhost:8088/profile/update", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(formData),
      });

      if (!response.ok) {
        throw new Error("Failed to update profile");
      }

      setIsEditing(false);
      fetchUserProfile(token!);
    } catch (error) {
      console.error(error);
    }
  };

  const handleClick = (id: number) => {
    navigate(`/products/${id}`);
  };

  const handleLogout = () => {
    localStorage.clear();
    navigate("/");
  };

  return (
      <div className="profile-container">
        <h1>User Profile</h1>
        {user ? (
            <div>
              <div className="profile-card">
                {!isEditing ? (
                    <div className="profile-details">
                      <p>
                        <strong>First Name:</strong> {user.first_name == null ? "" : user.first_name}
                      </p>
                      <p>
                        <strong>Last Name:</strong> {user.last_name == null ? "" : user.last_name}
                      </p>
                      <p>
                        <strong>Address:</strong> {user.address == null ? "" : user.address}
                      </p>
                      <p>
                        <strong>Phone Number:</strong> {user.phone_number == null ? "" : user.phone_number}
                      </p>
                      <p>
                        <strong>Email:</strong> {user.email}
                      </p>
                      <div className="profile-buttons">
                        <button onClick={() => setIsEditing(true)}>Edit Profile</button>
                        <button className="logout-button" onClick={handleLogout}>Logout</button>
                      </div>
                    </div>
                ) : (
                    <form onSubmit={handleSubmit} className="edit-form">
                      <label>
                        First Name:
                        <input
                            type="text"
                            name="firstName"
                            placeholder="Enter first name"
                            value={formData.first_name}
                            onChange={(e) =>
                                setFormData({ ...formData, first_name: e.target.value })
                            }
                        />
                      </label>
                      <label>
                        Last Name:
                        <input
                            type="text"
                            name="lastName"
                            placeholder="Enter last name"
                            value={formData.last_name}
                            onChange={(e) =>
                                setFormData({ ...formData, last_name: e.target.value })
                            }
                        />
                      </label>
                      <label>
                        Address:
                        <input
                            type="text"
                            name="address"
                            placeholder="Enter address"
                            value={formData.address}
                            onChange={(e) =>
                                setFormData({ ...formData, address: e.target.value })
                            }
                        />
                      </label>
                      <label>
                        Phone Number:
                        <input
                            type="text"
                            name="phoneNumber"
                            placeholder="Enter phone number"
                            value={formData.phone_number}
                            onChange={(e) =>
                                setFormData({ ...formData, phone_number: e.target.value })
                            }
                        />
                      </label>
                      <button type="submit">Save Changes</button>
                      <button type="button" onClick={() => setIsEditing(false)}>
                        Cancel
                      </button>
                    </form>
                )}
              </div>

              <div className="order-history">
                <h2>Order History</h2>
                {orders?.length > 0 ? (
                    <ul>
                      {orders.map((order, index) => (
                          <li
                              key={index}
                              className="order-item"
                              onClick={() => handleClick(order.id)}
                          >
                            <img
                                src={order.image}
                                alt={`${order.brand} ${order.model}`}
                            />
                            <div>
                              <p>
                                <strong>Type:</strong> {order.type}
                              </p>
                              <p>
                                <strong>Brand:</strong> {order.brand}
                              </p>
                              <p>
                                <strong>Model:</strong> {order.model}
                              </p>
                              <p>
                                <strong>Year:</strong> {order.year}
                              </p>
                              <p>
                                <strong>Total Price:</strong> ${order.total_price}
                              </p>
                              <p>
                                <strong>Start Date:</strong>{" "}
                                {new Date(order.start_date).toLocaleDateString()}
                              </p>
                              <p>
                                <strong>End Date:</strong>{" "}
                                {new Date(order.end_date).toLocaleDateString()}
                              </p>
                            </div>
                          </li>
                      ))}
                    </ul>
                ) : (
                    <p>No orders found</p>
                )}
              </div>
            </div>
        ) : (
            <p>Loading...</p>
        )}
      </div>
  );
};

export default UserProfile;