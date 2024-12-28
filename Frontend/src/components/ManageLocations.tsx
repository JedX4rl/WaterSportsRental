import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router";

interface Location {
  id: number;
  country: string;
  city: string;
  address: string;
  opening_time: string;
  closing_time: string;
  phone_number: string;
}

const LocationsPage: React.FC = () => {
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

  const [locations, setLocations] = useState<Location[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [isUpdateModalOpen, setIsUpdateModalOpen] = useState(false);
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [currentLocation, setCurrentLocation] = useState<Partial<Location>>({});
  const [showLocations, setShowLocations] = useState(false);

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

  const fetchLocations = async () => {
    try {
      const response = await fetch("http://localhost:8088/admin/locations", {
        headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
      });
      if (response.status === 401) {
        // Token expired or unauthorized
        navigate("/404");
        return;
      }
      if (!response.ok) {
        throw new Error("Failed to fetch locations");
      }
      const data = await response.json();
      setLocations(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    }
  };

  const handleDelete = async (id: number) => {
    try {
      const response = await fetch(
          `http://localhost:8088/admin/locations/?id=${id}`,
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
        throw new Error("Failed to delete location");
      }
      fetchLocations(); // Refresh the locations list
      setSuccess("Location deleted successfully");
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    }
  };

  const handleUpdate = async (id: number, updatedData: Partial<Location>) => {
    try {
      const response = await fetch(
          `http://localhost:8088/admin/locations/update`,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${localStorage.getItem("token")}`,
            },
            body: JSON.stringify({ id, ...updatedData }),
          },
      );
      if (response.status === 401) {
        // Token expired or unauthorized
        navigate("/404");
        return;
      }
      if (!response.ok) {
        throw new Error("Failed to update location");
      }
      fetchLocations(); // Refresh the locations list
      setIsUpdateModalOpen(false); // Close the update modal
      setSuccess("Location updated successfully");
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    }
  };

  const handleCreate = async (newLocation: Partial<Location>) => {
    try {
      const response = await fetch(`http://localhost:8088/admin/locations`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
        body: JSON.stringify(newLocation),
      });
      if (response.status === 401) {
        // Token expired or unauthorized
        navigate("/404");
        return;
      }
      if (!response.ok) {
        throw new Error("Failed to create location");
      }
      fetchLocations(); // Refresh the locations list
      setIsCreateModalOpen(false);
      setSuccess("Location created successfully");
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    }
  };

  const toggleShowLocations = () => {
    if (showLocations) {
      setShowLocations(false);
    } else {
      fetchLocations();
      setShowLocations(true);
    }
  };

  return (
      <div className="locations-page">
        <h1>Manage Locations</h1>
        <div className="button-container">
          <button onClick={toggleShowLocations}>
            {showLocations ? "Hide Locations" : "Get All Locations"}
          </button>
          <button onClick={() => setIsCreateModalOpen(true)}>
            Create New Location
          </button>
          <button onClick={() => navigate("/admin")}>AdminPage</button>
        </div>
        {error && <div className="error-message">{error}</div>}
        {success && <div className="success-message">{success}</div>}
        {showLocations && (
            <div className="locations-container">
              {locations.map((location) => (
                  <div className="location-card" key={location.id}>
                    <div className="card-content">
                      <button
                          className="delete-button"
                          onClick={() => handleDelete(location.id)}
                      >
                      </button>
                      <h3>
                        {location.country}, {location.city}
                      </h3>
                      <p>ID: {location.id}</p>
                      <p>Address: {location.address}</p>
                      <p>Opening Time: {location.opening_time}</p>
                      <p>Closing Time: {location.closing_time}</p>
                      <p>Phone: {location.phone_number}</p>
                      <button
                          className="update-button"
                          onClick={() => {
                            setCurrentLocation(location);
                            setIsUpdateModalOpen(true);
                          }}
                      >
                      </button>
                    </div>
                  </div>
              ))}
            </div>
        )}
        {isUpdateModalOpen && (
            <div className="modal">
              <div className="modal-content">
                <h2>Update Location</h2>
                <form
                    onSubmit={(e) => {
                      e.preventDefault();
                      handleUpdate(currentLocation.id!, currentLocation);
                    }}
                >
                  <input
                      type="text"
                      placeholder="Country"
                      value={currentLocation.country || ""}
                      onChange={(e) =>
                          setCurrentLocation({
                            ...currentLocation,
                            country: e.target.value,
                          })
                      }
                  />
                  <input
                      type="text"
                      placeholder="City"
                      value={currentLocation.city || ""}
                      onChange={(e) =>
                          setCurrentLocation({
                            ...currentLocation,
                            city: e.target.value,
                          })
                      }
                  />
                  <input
                      type="text"
                      placeholder="Address"
                      value={currentLocation.address || ""}
                      onChange={(e) =>
                          setCurrentLocation({
                            ...currentLocation,
                            address: e.target.value,
                          })
                      }
                  />
                  <input
                      type="text"
                      placeholder="Opening Time"
                      value={currentLocation.opening_time || ""}
                      onChange={(e) =>
                          setCurrentLocation({
                            ...currentLocation,
                            opening_time: e.target.value,
                          })
                      }
                  />
                  <input
                      type="text"
                      placeholder="Closing Time"
                      value={currentLocation.closing_time || ""}
                      onChange={(e) =>
                          setCurrentLocation({
                            ...currentLocation,
                            closing_time: e.target.value,
                          })
                      }
                  />
                  <input
                      type="text"
                      placeholder="Phone Number"
                      value={currentLocation.phone_number || ""}
                      onChange={(e) =>
                          setCurrentLocation({
                            ...currentLocation,
                            phone_number: e.target.value,
                          })
                      }
                  />
                  <button type="submit">Update</button>
                  <button type="button" onClick={() => setIsUpdateModalOpen(false)}>
                    Cancel
                  </button>
                </form>
              </div>
            </div>
        )}
        {isCreateModalOpen && (
            <div className="modal">
              <div className="modal-content">
                <h2>Create New Location</h2>
                <form
                    onSubmit={(e) => {
                      e.preventDefault();
                      handleCreate(currentLocation);
                    }}
                >
                  <input
                      type="text"
                      placeholder="Country"
                      value={currentLocation.country || ""}
                      onChange={(e) =>
                          setCurrentLocation({
                            ...currentLocation,
                            country: e.target.value,
                          })
                      }
                  />
                  <input
                      type="text"
                      placeholder="City"
                      value={currentLocation.city || ""}
                      onChange={(e) =>
                          setCurrentLocation({
                            ...currentLocation,
                            city: e.target.value,
                          })
                      }
                  />
                  <input
                      type="text"
                      placeholder="Address"
                      value={currentLocation.address || ""}
                      onChange={(e) =>
                          setCurrentLocation({
                            ...currentLocation,
                            address: e.target.value,
                          })
                      }
                  />
                  <input
                      type="text"
                      placeholder="Opening Time"
                      value={currentLocation.opening_time || ""}
                      onChange={(e) =>
                          setCurrentLocation({
                            ...currentLocation,
                            opening_time: e.target.value,
                          })
                      }
                  />
                  <input
                      type="text"
                      placeholder="Closing Time"
                      value={currentLocation.closing_time || ""}
                      onChange={(e) =>
                          setCurrentLocation({
                            ...currentLocation,
                            closing_time: e.target.value,
                          })
                      }
                  />
                  <input
                      type="text"
                      placeholder="Phone Number"
                      value={currentLocation.phone_number || ""}
                      onChange={(e) =>
                          setCurrentLocation({
                            ...currentLocation,
                            phone_number: e.target.value,
                          })
                      }
                  />
                  <button type="submit">Create</button>
                  <button type="button" onClick={() => setIsCreateModalOpen(false)}>
                    Cancel
                  </button>
                </form>
              </div>
            </div>
        )}
      </div>
  );
};

export default LocationsPage;