import React from "react";

function Option(props) {
  console.log("inside option", props);
  return (
    <button
      onClick={() => {
        props.onUpdate(props.option.id, {
          ...props.option,
          Votes: props.option.Votes + 1,
        });
      }}
    >
      {" "}
      <strong className="button-options">
        {" "}
        {props.option.Option} : {props.option.Votes}{" "}
      </strong>
    </button>
  );
}

export default Option;
