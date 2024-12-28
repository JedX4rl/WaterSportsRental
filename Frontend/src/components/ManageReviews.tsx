import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router";

interface Review {
  itemId: number;
  name: string;
  rating: number;
  comment: string;
  review_date: string;
  review_id: number;
}

const ReviewsPage: React.FC = () => {

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

  const [reviews, setReviews] = useState<Review[]>([]);
  const [expandedReview, setExpandedReview] = useState<number | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [showReviews, setShowReviews] = useState(false);


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

  const fetchReviews = async () => {
    try {
      const response = await fetch("http://localhost:8088/admin/reviews", {
        headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
      });
      if (response.status === 401) {
        navigate("/404");
        return;
      }
      if (!response.ok) {
        throw new Error("Failed to fetch reviews");
      }
      const data = await response.json();
      setReviews(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    }
  };

  const handleDelete = async (review_id: number) => {
    try {
      const response = await fetch(
        `http://localhost:8088/admin/reviews/?id=${review_id}`,
        {
          method: "DELETE",
          headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
        },
      );
      if (response.status === 401) {
        navigate("/404");
        return;
      }
      if (!response.ok) {
        throw new Error("Failed to delete review");
      }
      fetchReviews();
      setSuccess("Review deleted successfully");
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    }
  };

  const toggleShowReviews = () => {
    if (showReviews) {
      setShowReviews(false);
    } else {
      fetchReviews();
      setShowReviews(true);
    }
  };

  const toggleExpandReview = (review_id: number) => {
    setExpandedReview(expandedReview === review_id ? null : review_id);
  };

  return (
    <div className="reviews-page">
      <h1>Manage Reviews</h1>
      <div className="button-container">
        <button onClick={toggleShowReviews}>
          {showReviews ? "Hide Reviews" : "Get All Reviews"}
        </button>
        <button onClick={() => navigate("/admin")}>AdminPage</button>
      </div>
      {error && <div className="error-message">{error}</div>}
      {success && <div className="success-message">{success}</div>}
      {showReviews && (
        <div className="reviews-container">
          {reviews.map((review) => (
            <div
              className={`review-card ${expandedReview === review.review_id ? "expanded" : ""}`}
              key={review.review_id}
              onClick={() => toggleExpandReview(review.review_id)}
            >
              <div className="card-header">
                <h3>{review.name}</h3>
                <button
                  className="delete-button"
                  onClick={(e) => {
                    e.stopPropagation(); // Предотвращает раскрытие карточки при нажатии на кнопку
                    handleDelete(review.review_id);
                  }}
                >
                </button>
              </div>
              <div className="card-content">
                <p>Item ID: {review.itemId}</p>
                <p>Rating: {review.rating}</p>
                <p className="comment">
                  {expandedReview === review.review_id
                    ? review.comment
                    : `${review.comment.slice(0, 30)}...`}
                </p>
                <p>
                  Review Date:{" "}
                  {new Date(review.review_date).toLocaleDateString()}
                </p>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default ReviewsPage;
