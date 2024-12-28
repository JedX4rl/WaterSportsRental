import React, { useState, useEffect } from "react";

interface Location {
    id: number;
    country: string;
    city: string;
    address: string;
    opening_time: string;
    closing_time: string;
    phone_number: string;
}

const Contacts: React.FC = () => {
    const [locations, setLocations] = useState<Location[]>([]);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchLocations = async () => {
            try {
                const response = await fetch("http://localhost:8088/allLocations");
                if (!response.ok) {
                    throw new Error("Failed to fetch locations");
                }
                const data = await response.json();
                setLocations(data);
            } catch (err) {
                setError(err instanceof Error ? err.message : "An error occurred");
            }
        };

        fetchLocations();
    }, []);

    return (
        <div className="contacts-page">
            <h1>Store Locations</h1>
            {error && <div className="error-message">{error}</div>}
            <table className="locations-table">
                <thead>
                <tr>
                    <th>Country</th>
                    <th>City</th>
                    <th>Address</th>
                    <th>Opening Time</th>
                    <th>Closing Time</th>
                    <th>Phone Number</th>
                </tr>
                </thead>
                <tbody>
                {locations.map(location => (
                    <tr key={location.id}>
                        <td>{location.country}</td>
                        <td>{location.city}</td>
                        <td>{location.address}</td>
                        <td>{location.opening_time}</td>
                        <td>{location.closing_time}</td>
                        <td>{location.phone_number}</td>
                    </tr>
                ))}
                </tbody>
            </table>
        </div>
    );
};

export default Contacts;