import React, { useEffect } from "react";
import { useNavigate } from "react-router";

const AdminProfile: React.FC = () => {
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

  const handleNavigation = (path: string) => {
    navigate(path);
  };

  return (
    <div className="admin-profile-container">
      <h1>Admin Profile</h1>
      <div className="button-container">
        <button onClick={() => handleNavigation("/admin/locations")}>
          Locations
        </button>
        <button onClick={() => handleNavigation("/admin/users")}>Users</button>
        <button onClick={() => handleNavigation("/admin/reviews")}>
          Reviews
        </button>
        <button onClick={() => handleNavigation("/admin/items")}>Items</button>
        <button onClick={() => handleNavigation("/admin/dumps")}>Dumps</button>
      </div>
    </div>
  );
};

export default AdminProfile;
