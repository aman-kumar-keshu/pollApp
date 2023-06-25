import React from "react";
import { useNavigate } from "react-router-dom";

const Header = (props) => {
  const navigate = useNavigate();

  return (
    <>
      <div
        style={{
          position: "sticky",
          top: 0,
          padding: "10px 30px",
          height: "40px",
          background: "black",
          color: "white",
        }}
      >
        <div
          style={{
            display: "flex",
            justifyContent: "space-between",
          }}
        >
          <>
            <div>polling-app</div>{" "}
          </>
          <>
            <button
              onClick={() => {
                navigate("/new");
              }}
            >
              {" "}
              add new poll
            </button>
          </>
        </div>
      </div>
      <div className="card-header"></div>
      <h2 className="card-header-title header">
        You have {props.numPolls} Polls
      </h2>
    </>
  );
};
export default Header;
