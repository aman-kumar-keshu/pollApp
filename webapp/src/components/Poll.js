import React from "react";
import Option from "./Option";

function Poll(props) {
  return (
    <div className="poll">
      <ul>
        <div>
          {" "}
          <h2 className="title"> {props.poll.name} </h2>
        </div>
        <div>
          {" "}
          <h2 className="title"> TOPIC : {props.poll.topic}</h2>
        </div>
        <div>
          {" "}
          <img height={"200px"} width={"200px"} src={props.poll.src} />
        </div>

        {props.poll.options.map((option) => (
          <Option
            option={option}
            key={option.id}
            id={option.id}
            onUpdate={props.onUpdate}
            onDelete={props.onDelete}
          />
        ))}
        <div>
          <button
            onClick={() => {
              props.onDelete(props.poll.id);
            }}
          >
            <strong> Delete</strong>
          </button>
        </div>
      </ul>
    </div>
  );
}

export default Poll;
