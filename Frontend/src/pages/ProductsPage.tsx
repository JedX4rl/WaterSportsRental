import React, { useEffect, useState } from "react";
import { useSearchParams } from "react-router";
import Products from "@/components/Products.tsx";

const ProductsPage: React.FC = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  const [rentals, setRentals] = useState<never[]>([]);
  const [error, setError] = useState<string | null>(null);

  const fetchRentals = async (queryParams: URLSearchParams) => {
    console.log("fetchRentals called ");
    try {
      setError(null);
      console.log("fetchRentals trying to send info ");
      const response = await fetch(
        `http://localhost:8088/products/?${queryParams.toString()}`,
      );
      console.log("fetchRentals got response ");
      if (!response.ok) {
        throw new Error("Failed to fetch rentals");
      }
      const data = await response.json();
      setRentals(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
      setRentals([]); // Если ошибка, очищаем список товаров
    }
  };

  useEffect(() => {
    fetchRentals(searchParams);
  }, [searchParams]);

  const handleApplyFilters = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const formData = new FormData(e.currentTarget);
    const newParams: Record<string, string | string[]> = {};

    formData.forEach((value, key) => {
      if (value) {
        if (!newParams[key]) newParams[key] = [];
        (newParams[key] as string[]).push(value.toString());
      }
    });

    // Удаляем пустые значения для price-min и price-max
    if (!formData.get("price-min")) delete newParams["price-min"];
    if (!formData.get("price-max")) delete newParams["price-max"];

    const updatedSearchParams = new URLSearchParams();
    Object.entries(newParams).forEach(([key, value]) => {
      if (Array.isArray(value)) {
        value.forEach((val) => updatedSearchParams.append(key, val));
      } else {
        updatedSearchParams.append(key, value as string);
      }
    });

    setSearchParams(updatedSearchParams); // Обновляем URL
  };

  const locations = searchParams.getAll("location");
  const types = searchParams.getAll("type");
  const brands = searchParams.getAll("brand");

  return (
    <div className="products-page">
      <div className="container-products-page">
        {/* Фильтры отображаются всегда */}
        <div className="filters">
          <h3 className="filters-title">Filters</h3>
          <form className="filters-form" onSubmit={handleApplyFilters}>
            <div className="filter-group">
              <p>Location:</p>
              <label>
                <input
                  type="checkbox"
                  name="location"
                  value="1"
                  defaultChecked={locations.includes("bali")}
                />
                Bali
              </label>
              <label>
                <input
                  type="checkbox"
                  name="location"
                  value="2"
                  defaultChecked={locations.includes("hawaii")}
                />
                Hawaii
              </label>
              <label>
                <input
                  type="checkbox"
                  name="location"
                  value="3"
                  defaultChecked={locations.includes("sri-lanka")}
                />
                Sri Lanka
              </label>
            </div>
            <div className="filter-group">
              <p>Type:</p>
              <label>
                <input
                  type="checkbox"
                  name="type"
                  value="boat"
                  defaultChecked={types.includes("boats")}
                />
                Boats
              </label>
              <label>
                <input
                  type="checkbox"
                  name="type"
                  value="jet-ski"
                  defaultChecked={types.includes("jet-skis")}
                />
                Jet Skis
              </label>
              <label>
                <input
                  type="checkbox"
                  name="type"
                  value="paddleboards"
                  defaultChecked={types.includes("paddleboards")}
                />
                Paddleboards
              </label>
            </div>

            <div className="filter-group">
              <p>Brand:</p>
              <label>
                <input
                  type="checkbox"
                  name="brand"
                  value="kawasaki"
                  defaultChecked={brands.includes("kawasaki")}
                />
                Kawasaki
              </label>
              <label>
                <input
                  type="checkbox"
                  name="brand"
                  value="yamaha"
                  defaultChecked={brands.includes("yamaha")}
                />
                Yamaha
              </label>
            </div>

            <button type="submit" className="apply-filters-button">
              Apply Filters
            </button>
          </form>
        </div>

        {/* Основной контент */}
        <div className="products">
          {error ? (
            <div className="error-message">
              <p>Error: {error}</p>
            </div>
          ) : rentals == null || rentals.length === 0 ? (
            <div className="no-products">
              <p>No products match your filters.</p>
            </div>
          ) : (
            <Products rentals={rentals} />
          )}
        </div>
      </div>
    </div>
  );
};

export default ProductsPage;
