import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router";

interface User {
  id: number;
  email: string;
  first_name: string;
  last_name: string;
  address: string;
  phone_number: string;
  registration_date: string;
  role: string;
}

interface Agreement {
  date: string;
  item_id: number;
  start_date: string;
  end_date: string;
  status: boolean;
}

const UsersPage: React.FC = () => {
  const navigate = useNavigate();

  useEffect(() => {
    const checkPermission = async () => {
      try {
        const response = await fetch("http://localhost:8088/admin/permission", {
          headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
        });
        if (!response.ok) {
          navigate("/404");
        }
      } catch (error) {
        console.error("Error checking permissions:", error);
        navigate("/404");
      }
    };

    checkPermission();
  }, [navigate]);

  const [users, setUsers] = useState<User[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [showUsers, setShowUsers] = useState(false);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [currentUser, setCurrentUser] = useState<User | null>(null);
  const [agreements, setAgreements] = useState<Agreement[]>([]);

  useEffect(() => {
    if (error) {
      const timer = setTimeout(() => setError(null), 5000);
      return () => clearTimeout(timer);
    }
  }, [error]);

  useEffect(() => {
    if (success) {
      const timer = setTimeout(() => setSuccess(null), 5000);
      return () => clearTimeout(timer);
    }
  }, [success]);

  const fetchUsers = async () => {
    try {
      const response = await fetch("http://localhost:8088/admin/users/all", {
        headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
      });
      if (response.status === 401) {
        // Token expired or unauthorized
        navigate("*");
        return;
      }
      if (!response.ok) {
        throw new Error("Failed to fetch users");
      }
      const data = await response.json();
      setUsers(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    }
  };

  const fetchUserAgreements = async (id: number) => {
    try {
      const response = await fetch(`http://localhost:8088/admin/users/agreements?id=${id}`, {
        headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
      });
      if (!response.ok) {
        throw new Error("Failed to fetch user agreements");
      }
      const data = await response.json();
      setAgreements(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    }
  };

  const handleCardClick = async (user: User) => {
    setCurrentUser(user);
    await fetchUserAgreements(user.id);
    setIsModalOpen(true);
  };

  const handleDelete = async (id: number) => {
    try {
      const response = await fetch(
          `http://localhost:8088/admin/users/?id=${id}`,
          {
            method: "DELETE",
            headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
          },
      );
      if (response.status === 401) {
        // Token expired or unauthorized
        navigate("/404");
        return;
      }
      if (!response.ok) {
        throw new Error("Failed to delete user");
      }
      fetchUsers(); // Refresh the users list
      setSuccess("User deleted successfully");
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    }
  };

  const handleRoleChange = async (id: number, newRole: string) => {
    try {
      const response = await fetch(`http://localhost:8088/admin/users/role`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
        body: JSON.stringify({ id, role: newRole }),
      });
      if (response.status === 401) {
        // Token expired or unauthorized
        navigate("*");
        return;
      }
      if (!response.ok) {
        throw new Error("Failed to change user role");
      }
      fetchUsers(); // Refresh the users list
      setSuccess("User role updated successfully");
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    }
  };

  const toggleShowUsers = () => {
    if (showUsers) {
      setShowUsers(false);
    } else {
      fetchUsers();
      setShowUsers(true);
    }
  };

  return (
      <div className="users-page">
        <h1>Manage Users</h1>
        <div className="button-container">
          <button onClick={toggleShowUsers}>
            {showUsers ? "Hide Users" : "Get All Users"}
          </button>
          <button onClick={() => navigate("/admin")}>AdminPage</button>
        </div>
        {error && <div className="error-message">{error}</div>}
        {success && <div className="success-message">{success}</div>}
        {showUsers && (
            <div className="users-container">
              {users.map((user) => (
                  <div className="user-card" key={user.id} onClick={() => handleCardClick(user)}>
                    <button
                        className="delete-button"
                        onClick={(e) => {
                          e.stopPropagation();
                          handleDelete(user.id);
                        }}
                    >
                    </button>
                    <div className="card-content">
                      <h3>
                        {user.first_name} {user.last_name}
                      </h3>
                      <p>ID: {user.id}</p>
                      <p>Email: {user.email}</p>
                      <p>Address: {user.address}</p>
                      <p>Phone: {user.phone_number}</p>
                      <p>
                        Registration Date:{" "}
                        {new Date(user.registration_date).toLocaleDateString()}
                      </p>
                      <p>Role: {user.role}</p>
                      <div className="role-change">
                        <label>Change Role:</label>
                        <select
                            value={user.role}
                            onChange={(e) => handleRoleChange(user.id, e.target.value)}
                        >
                          <option value="user">User</option>
                          <option value="admin">Admin</option>
                        </select>
                      </div>
                    </div>
                  </div>
              ))}
            </div>
        )}
        {isModalOpen && currentUser && (
            <div className="modal">
              <div className="modal-content">
                <h2>Agreements for {currentUser.first_name} {currentUser.last_name}</h2>
                <table>
                  <thead>
                  <tr>
                    <th>Date</th>
                    <th>Item ID</th>
                    <th>Start Date</th>
                    <th>End Date</th>
                    <th>Status</th>
                  </tr>
                  </thead>
                  <tbody>
                  {agreements?.map((agreement, index) => (
                      <tr key={index}>
                        <td>{new Date(agreement.date).toLocaleString()}</td>
                        <td>{agreement.item_id}</td>
                        <td>{new Date(agreement.start_date).toLocaleString()}</td>
                        <td>{new Date(agreement.end_date).toLocaleString()}</td>
                        <td>{agreement.status ? "Active" : "Inactive"}</td>
                      </tr>
                  ))}
                  </tbody>
                </table>
                <button onClick={() => setIsModalOpen(false)}>Close</button>
              </div>
            </div>
        )}
      </div>
  );
};

export default UsersPage;