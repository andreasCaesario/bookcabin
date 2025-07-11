
import React, { useState, useEffect, useRef } from "react";
import "./App.css";

function App() {
  const [flightNumber, setFlightNumber] = useState("");
  const [date, setDate] = useState("");
  const [name, setName] = useState("");
  const [crewId, setCrewId] = useState("");
  const [aircraft, setAircraft] = useState("");
  const [aircraftList, setAircraftList] = useState([]);
  const [checkResult, setCheckResult] = useState(null);
  const [generateResult, setGenerateResult] = useState(null);
  const [error, setError] = useState("");
  const [showPopup, setShowPopup] = useState(false);
  const [popupType, setPopupType] = useState(""); // 'success' or 'error'
  const popupTimeout = useRef(null);
  // Helper to get current date in DD - MM - YY format
  function getCurrentDatePlaceholder() {
    const now = new Date();
    const dd = String(now.getDate()).padStart(2, '0');
    const mm = String(now.getMonth() + 1).padStart(2, '0');
    const yy = String(now.getFullYear()).slice(-2);
    return `${dd} - ${mm} - ${yy}`;
  }

  useEffect(() => {
    fetch("/api/aircraft-list")
      .then(res => res.json())
      .then(data => setAircraftList(data.map(item => item.type)))
      .catch(() => setAircraftList([]));
  }, []);


  const handleGenerate = async () => {
    setError("");
    setGenerateResult(null);
    setShowPopup(false);
    setPopupType("");
    // Simple input validation
    if (!flightNumber.trim() || !date.trim() || !name.trim() || !crewId.trim() || !aircraft.trim()) {
      setError("Please fill in all fields before generating a voucher.");
      setPopupType("error");
      setShowPopup(true);
      if (popupTimeout.current) clearTimeout(popupTimeout.current);
      popupTimeout.current = setTimeout(() => setShowPopup(false), 5000);
      return;
    }
    try {
      const checkRes = await fetch("/api/check", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ flightNumber, date })
      });
      const checkData = await checkRes.json();
      if (checkData.error) {
        setError(checkData.error);
        setPopupType("error");
        setShowPopup(true);
        if (popupTimeout.current) clearTimeout(popupTimeout.current);
        popupTimeout.current = setTimeout(() => setShowPopup(false), 7000);
        return;
      }
      if (checkData.exists) {
        setError("Sorry, voucher already generated for this flight number and for the date you choose.");
        setPopupType("error");
        setShowPopup(true);
        if (popupTimeout.current) clearTimeout(popupTimeout.current);
        popupTimeout.current = setTimeout(() => setShowPopup(false), 7000);
        return;
      }
      // If not exists, generate
      const res = await fetch("/api/generate", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name, id: crewId, flightNumber, date, aircraft })
      });
      const data = await res.json();
      setGenerateResult(data);
      if (!data.success) {
        setError(data.error || "Failed to generate voucher");
        setPopupType("error");
        setShowPopup(true);
        if (popupTimeout.current) clearTimeout(popupTimeout.current);
        popupTimeout.current = setTimeout(() => setShowPopup(false), 7000);
      } else {
        setPopupType("success");
        setShowPopup(true);
        if (popupTimeout.current) clearTimeout(popupTimeout.current);
        popupTimeout.current = setTimeout(() => setShowPopup(false), 5000);
      }
    } catch (e) {
      setError("Generate failed");
      setPopupType("error");
      setShowPopup(true);
      if (popupTimeout.current) clearTimeout(popupTimeout.current);
      popupTimeout.current = setTimeout(() => setShowPopup(false), 7000);
    }
  };

  return (
    <div className="app-container">
      <h2>Voucher Seat Assignment</h2>
      <div className="form-group">
        <label>Flight Number:</label>
        <input value={flightNumber} onChange={e => setFlightNumber(e.target.value)} type="text" />
      </div>
      <div className="form-group">
        <label>Date (DD - MM - YY):</label>
        <input
          type="text"
          value={date}
          onChange={e => {
            // Only allow numbers and dashes, and auto-format
            let v = e.target.value.replace(/[^0-9-]/g, "");
            if (v.length === 2 || v.length === 5) {
              if (date.length < v.length) v += " - ";
            }
            setDate(v);
          }}
          placeholder={getCurrentDatePlaceholder()}
          maxLength={12}
        />
      </div>
      <div className="form-group">
        <label>Name:</label>
        <input value={name} onChange={e => setName(e.target.value)} type="text" />
      </div>
      <div className="form-group">
        <label>Crew ID:</label>
        <input value={crewId} onChange={e => setCrewId(e.target.value)} type="text" />
      </div>
      <div className="form-group">
        <label>Aircraft:</label>
        <select value={aircraft} onChange={e => setAircraft(e.target.value)}>
          <option value="">-- Aircraft type --</option>
          {aircraftList.map(type => (
            <option key={type} value={type}>{type}</option>
          ))}
        </select>
      </div>
      <div className="button-row">
        <button onClick={handleGenerate}>Generate Vouchers</button>
      </div>
      {/* Popup for error/success */}
      {showPopup && (
        <div className={`popup-message ${popupType}`}>
          {popupType === "success" && generateResult && generateResult.success && (
            <>
              <strong>Success!</strong> Seats: {generateResult.seats.join(", ")}
            </>
          )}
          {popupType === "error" && error && (
            <>{error}</>
          )}
        </div>
      )}
    </div>
  );
}

export default App;
