import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router";

interface Dump {
    filename: string;
}

interface Log {
    id: number;
    itemId: number;
    startDate: string;
    endDate: string;
}

const DumpsPage: React.FC = () => {

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

    const [dumps, setDumps] = useState<Dump[]>([]);
    const [logs, setLogs] = useState<Log[]>([]);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);
    const [showDumps, setShowDumps] = useState(false);
    const [showLogs, setShowLogs] = useState(false);

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

    const fetchDumps = async () => {
        try {
            const token = localStorage.getItem("token");
            const response = await fetch("http://localhost:8088/admin/database/all", {
                headers: { Authorization: `Bearer ${token}` },
            });
            if (!response.ok) {
                throw new Error("Failed to fetch dumps");
            }
            const data = await response.json();
            setDumps(data);
        } catch (err) {
            setError(err instanceof Error ? err.message : "An error occurred");
        }
    };

    const fetchLogs = async () => {
        try {
            const token = localStorage.getItem("token");
            const response = await fetch("http://localhost:8088/admin/logs", {
                headers: { Authorization: `Bearer ${token}` },
            });
            if (!response.ok) {
                throw new Error("Failed to fetch logs");
            }
            const data = await response.json();
            setLogs(data);
            setShowLogs(true);
        } catch (err) {
            setError(err instanceof Error ? err.message : "An error occurred");
        }
    };

    const restoreDump = async (filename: string) => {
        try {
            const token = localStorage.getItem("token");
            const response = await fetch("http://localhost:8088/admin/database/restore", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    Authorization: `Bearer ${token}`
                },
                body: JSON.stringify({ filename })
            });
            if (!response.ok) {
                throw new Error("Failed to restore dump");
            }
            setSuccess("Restore successful");
        } catch (err) {
            setError(err instanceof Error ? err.message : "An error occurred");
        }
    };

    const createDump = async () => {
        try {
            const token = localStorage.getItem("token");
            const response = await fetch("http://localhost:8088/admin/database/dump", {
                method: "GET",
                headers: { Authorization: `Bearer ${token}` },
            });
            if (!response.ok) {
                throw new Error("Failed to create dump");
            }
            setSuccess("Dump created successfully");
            fetchDumps(); // Обновляем список дампов после создания нового
        } catch (err) {
            setError(err instanceof Error ? err.message : "An error occurred");
        }
    };

    const handleRestore = (filename: string) => {
        if (window.confirm(`Are you sure you want to restore from ${filename}?`)) {
            restoreDump(filename);
        }
    };

    const toggleShowDumps = () => {
        if (showDumps) {
            setShowDumps(false);
        } else {
            fetchDumps();
            setShowDumps(true);
        }
    };

    const toggleShowLogs = () => {
        if (showLogs) {
            setShowLogs(false);
        } else {
            fetchLogs();
        }
    };

    return (
        <div className="dumps-page">
            <h1>Manage Dumps</h1>
            <div className="button-container">
                <button onClick={toggleShowDumps}>
                    {showDumps ? "Hide Dumps" : "Get All Dumps"}
                </button>
                <button onClick={createDump}>Create Dump</button>
                <button onClick={() => navigate("/admin")}>AdminPage</button>
                <button onClick={toggleShowLogs}>
                    {showLogs ? "Hide Logs" : "Get Logs"}
                </button>
            </div>
            {error && <div className="error-message">{error}</div>}
            {success && <div className="success-message">{success}</div>}
            {showDumps && (
                <table className="dumps-table">
                    <thead>
                    <tr>
                        <th>Filename</th>
                    </tr>
                    </thead>
                    <tbody>
                    {dumps.map((dump) => (
                        <tr key={dump.filename} onClick={() => handleRestore(dump.filename)}>
                            <td>{dump.filename}</td>
                        </tr>
                    ))}
                    </tbody>
                </table>
            )}
            {showLogs && (
                <table className="logs-table">
                    <thead>
                    <tr>
                        <th>ID</th>
                        <th>Item ID</th>
                        <th>Start Date</th>
                        <th>End Date</th>
                    </tr>
                    </thead>
                    <tbody>
                    {logs.map((log) => (
                        <tr key={log.id}>
                            <td>{log.id}</td>
                            <td>{log.itemId}</td>
                            <td>{log.startDate}</td>
                            <td>{log.endDate}</td>
                        </tr>
                    ))}
                    </tbody>
                </table>
            )}
        </div>
    );
};

export default DumpsPage;