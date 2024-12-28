import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router";

interface Item {
  id: number;
  type: string;
  brand: string;
  model: string;
  year: number;
  price: number;
  image: string;
}

interface Location {
  id: number;
  address: string;
}

const ItemsPage: React.FC = () => {
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

  const [items, setItems] = useState<Item[]>([]);
  const [locations, setLocations] = useState<Location[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [isUpdateModalOpen, setIsUpdateModalOpen] = useState(false);
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [currentItem, setCurrentItem] = useState<Partial<Item>>({});
  const [selectedLocation, setSelectedLocation] = useState<number | null>(null);
  const [quantity, setQuantity] = useState<number>(1);
  const [showItems, setShowItems] = useState(false);

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

  const fetchItems = async () => {
    try {
      const response = await fetch("http://localhost:8088/admin/items", {
        headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
      });
      if (response.status === 401) {
        // Token expired or unauthorized
        navigate("*");
        return;
      }
      if (!response.ok) {
        throw new Error("Failed to fetch items");
      }
      const data = await response.json();
      console.log(data);
      setItems(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    }
  };

  const fetchLocations = async () => {
    try {
      const response = await fetch("http://localhost:8088/admin/locations", {
        headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
      });
      if (!response.ok) {
        throw new Error("Failed to fetch locations");
      }
      const data = await response.json();
      setLocations(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    }
  };

  const handleCardClick = async (item: Item) => {
    setCurrentItem(item);
    await fetchLocations();
    setIsAddModalOpen(true);
  };

  const handleDelete = async (id: number) => {
    try {
      const response = await fetch(
          `http://localhost:8088/admin/items/?id=${id}`,
          {
            method: "DELETE",
            headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
          },
      );
      if (response.status === 401) {
        // Token expired or unauthorized
        navigate("*");
        return;
      }
      if (!response.ok) {
        throw new Error("Failed to delete item");
      }
      fetchItems(); // Refresh the items list
      setSuccess("Item deleted successfully");
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    }
  };

  const handleUpdate = async (id: number, updatedData: Partial<Item>) => {
    try {
      const response = await fetch(`http://localhost:8088/admin/items/update`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
        body: JSON.stringify({ id, ...updatedData }),
      });
      if (response.status === 401) {
        // Token expired or unauthorized
        navigate("*");
        return;
      }
      if (!response.ok) {
        throw new Error("Failed to update item");
      }
      fetchItems(); // Refresh the items list
      setIsUpdateModalOpen(false); // Close the update modal
      setSuccess("Item updated successfully");
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    }
  };

  const handleCreate = async (newItem: Partial<Item>) => {
    try {
      const response = await fetch(`http://localhost:8088/admin/items`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
        body: JSON.stringify(newItem),
      });
      if (response.status === 401) {
        // Token expired or unauthorized
        navigate("*");
        return;
      }
      if (!response.ok) {
        throw new Error("Failed to create item");
      }
      fetchItems(); // Refresh the items list
      setIsCreateModalOpen(false); // Close the create modal
      setSuccess("Item created successfully");
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    }
  };

  const handleAddItem = async (e: React.FormEvent) => {
    e.preventDefault();
    if (currentItem.id && selectedLocation && quantity > 0) {
      try {
        const response = await fetch("http://localhost:8088/admin/items/add", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${localStorage.getItem("token")}`,
          },
          body: JSON.stringify({
            product_id: currentItem.id,
            location_id: selectedLocation,
            number: quantity,
          }),
        });
        if (!response.ok) {
          throw new Error("Failed to add item");
        }
        setSuccess("Item added successfully");
        setIsAddModalOpen(false);
      } catch (err) {
        setError(err instanceof Error ? err.message : "An error occurred");
      }
    }
  };

  const toggleShowItems = () => {
    if (showItems) {
      setShowItems(false);
    } else {
      fetchItems();
      setShowItems(true);
    }
  };

  return (
      <div className="items-page">
        <h1>Manage Items</h1>
        <div className="button-container">
          <button onClick={toggleShowItems}>
            {showItems ? "Hide Items" : "Get All Items"}
          </button>
          <button onClick={() => setIsCreateModalOpen(true)}>
            Create New Item
          </button>
          <button onClick={() => navigate("/admin")}>AdminPage</button>
        </div>
        {error && <div className="error-message">{error}</div>}
        {success && <div className="success-message">{success}</div>}
        {showItems && (
            <div className="items-container">
              {items.map((item) => (
                  <div className="product-card" key={item.id} onClick={() => handleCardClick(item)}>
                    <button
                        className="delete-button"
                        onClick={(e) => {
                          e.stopPropagation();
                          handleDelete(item.id);
                        }}
                    >
                    </button>
                    <div className="card-image">
                      <img src={item.image} alt={`${item.brand} ${item.model}`} />
                    </div>
                    <div className="card-content">
                      <h3>
                        {item.brand} {item.model}
                      </h3>
                      <p>ID: {item.id}</p>
                      <p>Type: {item.type}</p>
                      <p>Year: {item.year}</p>
                      <p>
                        Price: <strong>${item.price}/hour</strong>
                      </p>
                      <p>Image: {item.image}</p>
                    </div>
                    <button
                        className="update-button"
                        onClick={(e) => {
                          e.stopPropagation();
                          setCurrentItem(item);
                          setIsUpdateModalOpen(true);
                        }}
                    >
                    </button>
                  </div>
              ))}
            </div>
        )}
        {isUpdateModalOpen && (
            <div className="modal">
              <div className="modal-content">
                <h2>Update Item</h2>
                <form
                    onSubmit={(e) => {
                      e.preventDefault();
                      handleUpdate(currentItem.id!, currentItem);
                    }}
                >
                  <input
                      type="text"
                      placeholder="Type"
                      value={currentItem.type || ""}
                      onChange={(e) =>
                          setCurrentItem({ ...currentItem, type: e.target.value })
                      }
                  />
                  <input
                      type="text"
                      placeholder="Brand"
                      value={currentItem.brand || ""}
                      onChange={(e) =>
                          setCurrentItem({ ...currentItem, brand: e.target.value })
                      }
                  />
                  <input
                      type="text"
                      placeholder="Model"
                      value={currentItem.model || ""}
                      onChange={(e) =>
                          setCurrentItem({ ...currentItem, model: e.target.value })
                      }
                  />
                  <input
                      type="number"
                      placeholder="Year"
                      value={currentItem.year || ""}
                      onChange={(e) =>
                          setCurrentItem({
                            ...currentItem,
                            year: parseInt(e.target.value),
                          })
                      }
                  />
                  <input
                      type="number"
                      placeholder="Price"
                      value={currentItem.price || ""}
                      onChange={(e) =>
                          setCurrentItem({
                            ...currentItem,
                            price: parseFloat(e.target.value),
                          })
                      }
                  />
                  <input
                      type="text"
                      placeholder="Image URL"
                      value={currentItem.image || ""}
                      onChange={(e) =>
                          setCurrentItem({ ...currentItem, image: e.target.value })
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
                <h2>Create New Item</h2>
                <form
                    onSubmit={(e) => {
                      e.preventDefault();
                      handleCreate(currentItem);
                    }}
                >
                  <input
                      type="text"
                      placeholder="Type"
                      value={currentItem.type || ""}
                      onChange={(e) =>
                          setCurrentItem({ ...currentItem, type: e.target.value })
                      }
                  />
                  <input
                      type="text"
                      placeholder="Brand"
                      value={currentItem.brand || ""}
                      onChange={(e) =>
                          setCurrentItem({ ...currentItem, brand: e.target.value })
                      }
                  />
                  <input
                      type="text"
                      placeholder="Model"
                      value={currentItem.model || ""}
                      onChange={(e) =>
                          setCurrentItem({ ...currentItem, model: e.target.value })
                      }
                  />
                  <input
                      type="number"
                      placeholder="Year"
                      value={currentItem.year || ""}
                      onChange={(e) =>
                          setCurrentItem({
                            ...currentItem,
                            year: parseInt(e.target.value),
                          })
                      }
                  />
                  <input
                      type="number"
                      placeholder="Price"
                      value={currentItem.price || ""}
                      onChange={(e) =>
                          setCurrentItem({
                            ...currentItem,
                            price: parseFloat(e.target.value),
                          })
                      }
                  />
                  <input
                      type="text"
                      placeholder="Image URL"
                      value={currentItem.image || ""}
                      onChange={(e) =>
                          setCurrentItem({ ...currentItem, image: e.target.value })
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
        {isAddModalOpen && (
            <div className="modal">
              <div className="modal-content">
                <h2>Add Item to Location</h2>
                <form onSubmit={handleAddItem}>
                  <label>
                    Location:
                    <select
                        value={selectedLocation ?? ""}
                        onChange={(e) => setSelectedLocation(parseInt(e.target.value))}
                        required
                    >
                      <option value="" disabled>
                        Select location
                      </option>
                      {locations.map((location) => (
                          <option key={location.id} value={location.id}>
                            {location.address}
                          </option>
                      ))}
                    </select>
                  </label>
                  <label>
                    Quantity:
                    <input
                        type="number"
                        value={quantity}
                        onChange={(e) => setQuantity(parseInt(e.target.value))}
                        min="1"
                        required
                    />
                  </label>
                  <button type="submit">Submit</button>
                  <button type="button" onClick={() => setIsAddModalOpen(false)}>
                    Cancel
                  </button>
                </form>
              </div>
            </div>
        )}
      </div>
  );
};

export default ItemsPage;