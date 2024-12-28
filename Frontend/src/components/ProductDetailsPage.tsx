
import React, { useEffect, useState } from "react";
import { useParams } from "react-router";
import DatePicker from "react-datepicker";
import "react-datepicker/dist/react-datepicker.css";

interface ProductDetails {
  image: string;
  type: string;
  brand: string;
  model: string;
  year: number;
  price: number;
}

interface Review {
  itemId?: number;
  name: string;
  comment: string;
  rating: number;
  review_date: string;
}

interface Location {
  product_id: string;
  location_id: string;
  location: string;
  number: string;
}

interface AvailableDate {
  start_date: string;
  end_date: string;
}

const ProductDetailsPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const parsedId: number | undefined = id ? +id : undefined;
  const [product, setProduct] = useState<ProductDetails | null>(null);
  const [reviews, setReviews] = useState<Review[]>([]);
  const [locations, setLocations] = useState<Location[]>([]);
  const [availableDates, setAvailableDates] = useState<AvailableDate[]>([]);
  const [highlightedDates, setHighlightedDates] = useState<Date[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [expandedReview, setExpandedReview] = useState<number | null>(null);
  const [newReview, setNewReview] = useState({
    itemId: parsedId,
    comment: "",
    rating: "",
  });
  const [showBookingForm, setShowBookingForm] = useState(false);
  const [bookingData, setBookingData] = useState({
    location_id: "",
    start_date: new Date(),
    end_date: new Date(),
    payment_method: "card",
  });
  const [message, setMessage] = useState<{ text: string; type: "success" | "error" } | null>(null);

  const fetchProductDetails = async () => {
    try {
      const response = await fetch(`http://localhost:8088/products?id=${id}`);
      if (!response.ok) throw new Error("Failed to fetch product details");
      const data = await response.json();
      setProduct(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    }
  };

  const fetchReviews = async () => {
    try {
      const response = await fetch(`http://localhost:8088/products/reviews?id=${id}`);
      if (!response.ok) throw new Error("Failed to fetch reviews");
      const data = await response.json();
      setReviews(data);
    } catch (err) {
      console.error("Failed to fetch reviews", err);
    }
  };

  const fetchLocations = async () => {
    try {
      const response = await fetch(`http://localhost:8088/products/locations?id=${id}`);
      if (!response.ok) throw new Error("Failed to fetch locations");
      const data: Location[] = await response.json();
      setLocations(data);
    } catch (err) {
      console.error("Failed to fetch locations", err);
    }
  };

  const fetchAvailableDates = async (locationId: number, itemId: number) => {
    try {
      const response = await fetch(`http://localhost:8088/products/dates?location_id=${locationId}&item_id=${itemId}`);
      if (!response.ok) throw new Error("Failed to fetch available dates");
      const data = await response.json();
      console.log("Fetched available dates:", data);
      setAvailableDates(data);
      console.log(availableDates)
    } catch (err) {
      console.error("Failed to fetch available dates", err);
      setAvailableDates([]);
    }
  };

  useEffect(() => {
    fetchProductDetails();
    fetchReviews();
    fetchLocations();
  }, [id]);

  useEffect(() => {
    console.log(availableDates)
    const highlightWithRanges = availableDates.flatMap((period) => {
      const startDate = new Date(period.start_date);
      const endDate = new Date(period.end_date);
      const dates = [];
      for (let d = new Date(startDate); d <= endDate; d.setDate(d.getDate() + 1)) {
        dates.push(new Date(d));
      }
      console.log("ADADAD")
      console.log(dates)
      return dates;
    });
    console.log(highlightWithRanges)
    setHighlightedDates(highlightWithRanges);
  }, [availableDates]);

  const handleLocationChange = async (e: React.ChangeEvent<HTMLSelectElement>) => {
    const location_id = e.target.value;
    setBookingData((prevData) => ({
      ...prevData,
      location_id,
    }));

    if (parsedId) {
      await fetchAvailableDates(+location_id, parsedId);
    }
  };

  const handleShowBookingForm = () => {
    setShowBookingForm(true);
  };

  const handleSubmitReview = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newReview.comment || +newReview.rating <= 0) {
      alert("Please fill all fields correctly.");
      return;
    }

    const newReviewInt = {
      itemId: newReview.itemId,
      comment: newReview.comment,
      rating: +newReview.rating,
    };

    const token = localStorage.getItem("token");
    try {
      const response = await fetch(`http://localhost:8088/products/review`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(newReviewInt),
      });

      if (!response.ok) throw new Error("Failed to submit review");

      fetchReviews();

      setNewReview({
        itemId: parsedId,
        comment: "",
        rating: "",
      });
      setMessage({ text: "Review submitted successfully!", type: "success" });
    } catch (err) {
      setMessage({ text: "Failed to submit review.", type: "error" });
      console.error("Error submitting review", err);
    } finally {
      setTimeout(() => setMessage(null), 5000);
    }
  };

  const handleSubmitBooking = async (e: React.FormEvent) => {
    e.preventDefault();
    const token = localStorage.getItem("token");
    const bookingRequest = {
      ids: [parsedId!],
      location_id: +bookingData.location_id,
      start_date: bookingData.start_date,
      end_date: bookingData.end_date,
      payment_method: bookingData.payment_method,
    };

    try {
      const response = await fetch(`http://localhost:8088/products/rent`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify(bookingRequest),
      });

      if (!response.ok) throw new Error("Failed to book product");

      setMessage({ text: "Booking successful!", type: "success" });
      setShowBookingForm(false);
      fetchLocations(); // Обновляем количество доступных товаров после бронирования
      setBookingData({
        location_id: "",
        start_date: new Date(),
        end_date: new Date(),
        payment_method: "card",
      });
    } catch (err) {
      setMessage({ text: "Booking failed. Please try again.", type: "error" });
      console.error("Error booking product", err);
    } finally {
      setTimeout(() => setMessage(null), 5000);
    }
  };

  const handleStartDateChange = (date: Date | null) => {
    if (date) {
      setBookingData({ ...bookingData, start_date: date });
    }
  };

  const handleEndDateChange = (date: Date | null) => {
    if (date) {
      setBookingData({ ...bookingData, end_date: date });
    }
  };

  if (error) {
    return <div className="error-message">Error: {error}</div>;
  }

  if (!product) {
    return <div className="loading-message">Loading...</div>;
  }

  const toggleExpandReview = (review_id: number) => {
    setExpandedReview(expandedReview === review_id ? null : review_id);
  };

  return (
      <div className="product-details-container">
        {message && (
            <div className={`message ${message.type}`}>
              {message.text}
            </div>
        )}
        <div className="product-info">
          <img
              className="product-image"
              src={product.image?.startsWith("https") ? product.image : "/" + product.image}
              alt={`${product.brand} ${product.model}`}
          />
          <div className="product-text-info">
            <h1>
              {product.brand} {product.model}
            </h1>
            <p>
              <strong>Type:</strong> {product.type}
            </p>
            <p>
              <strong>Year:</strong> {product.year}
            </p>
            <p>
              <strong>Price:</strong> ${product.price}/hour
            </p>
            <h3>Available at:</h3>
            {locations?.length > 0 ? (
                <ul className="locations-list">
                  {locations?.map((location, index) => (
                      <li key={index}>
                        {location.location} ({location.number} available)
                      </li>
                  ))}
                </ul>
            ) : (
                <p>No locations available.</p>
            )}
            <button onClick={handleShowBookingForm}>
              Book Product
            </button>
          </div>
        </div>
        {showBookingForm && (
            <form className="booking-form" onSubmit={handleSubmitBooking}>
              <h3>Book Product</h3>
              <label>
                Location:
                <select
                    value={bookingData.location_id}
                    onChange={handleLocationChange}
                    required
                >
                  <option value="" disabled>
                    Select location
                  </option>
                  {locations?.map((location) => (
                      <option key={location.location_id} value={location.location_id}>
                        {location.location}
                      </option>
                  ))}
                </select>
              </label>
              <label>
                Start Date:
                <DatePicker
                    selected={bookingData.start_date}
                    onChange={handleStartDateChange}
                    dateFormat="yyyy/MM/dd"
                    required
                    highlightDates={highlightedDates}
                    dayClassName={(date) => {
                      return highlightedDates.some(
                          (highlightDate) => highlightDate.getTime() === date.getTime()
                      )
                          ? "react-datepicker__day--highlighted-custom-available"
                          : "";
                    }}
                />
              </label>
              <label>
                End Date:
                <DatePicker
                    selected={bookingData.end_date}
                    onChange={handleEndDateChange}
                    dateFormat="yyyy/MM/dd"
                    required
                    highlightDates={highlightedDates}
                    dayClassName={(date) => {
                      return highlightedDates.some(
                          (highlightDate) => highlightDate.getTime() === date.getTime()
                      )
                          ? "react-datepicker__day--highlighted-custom-available"
                          : "";
                    }}
                />
              </label>
              <label>
                Payment Method:
                <select
                    value={bookingData.payment_method}
                    onChange={(e) =>
                        setBookingData({
                          ...bookingData,
                          payment_method: e.target.value,
                        })
                    }
                    required
                >
                  <option value="card">Card</option>
                  <option value="cash">Cash</option>
                </select>
              </label>
              <button type="submit">Submit Booking</button>
              <button type="button" onClick={() => setShowBookingForm(false)}>
                Cancel
              </button>
            </form>
        )}
        <div className="product-reviews">
          <h3>Reviews:</h3>
          {reviews?.length > 0 ? (
              reviews?.map((review, index) => (
                  <div
                      key={index}
                      className={`item-review-card ${
                          expandedReview === index ? "expanded" : ""
                      }`}
                      onClick={() => toggleExpandReview(index)}
                  >
                    <div className="card-content">
                      <p>Rating: {review.rating}/5</p>
                      <p>
                        Review Date:{" "}
                        {new Date(review.review_date).toLocaleDateString()}
                      </p>
                      <p>
                        <strong>{review.name}</strong>
                      </p>
                      <p className="comment">
                        {expandedReview === index
                            ? review.comment
                            : `${review.comment.slice(0, 30)}...`}
                      </p>
                    </div>
                  </div>
              ))
          ) : (
              <p>No reviews yet.</p>
          )}
          <h3>Leave a Review</h3>
          <form className="review-form" onSubmit={handleSubmitReview}>
          <textarea
              placeholder="Your Review (max 400 characters)"
              value={newReview.comment}
              maxLength={400}
              onChange={(e) =>
                  setNewReview({
                    ...newReview,
                    comment: e.target.value.slice(0, 400),
                  })
              }
              required
          />
            <input
                type="number"
                placeholder="Rating (1-5)"
                min="1"
                max="5"
                value={newReview.rating}
                onChange={(e) =>
                    setNewReview({ ...newReview, rating: e.target.value })
                }
                required
            />
            <button type="submit">Submit</button>
          </form>
        </div>
      </div>
  );
};

export default ProductDetailsPage;