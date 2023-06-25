import React, { useEffect, useState } from "react";
import axios from "axios";
import PollList from "../components/PollList";
import Header from "../components/Header";
const endpoint = `http://localhost:8080`;

const sortPolls = (polls) => polls.sort((poll1, poll2) => poll1.id - poll2.id);

function Home() {
  const [polls, setPolls] = useState([]);

  useEffect(() => {
    performAPICall();
  }, []);

  const performAPICall = async () => {
    console.log("Making call to backend to fetch polls");
    await axios.get(`${endpoint}/polls`).then((res) => {
      setPolls(sortPolls(res.data.polls.items));
    });
  };

  const handleUpdate = async (id, updatedPoll) => {
    console.log("Updating the poll", updatedPoll);

    const oldListWithoutId = polls.filter((poll) => poll.id !== id);
    console.log("oldListWithoutId", oldListWithoutId);

    const newPolls = sortPolls([...oldListWithoutId, updatedPoll]);
    console.log("newPolls", newPolls);
    setPolls(newPolls);
    await axios.put(`${endpoint}/poll/${id}`, updatedPoll);
  };

  const handleDelete = (id) => {
    console.log("handle delete");
  };

  return (
    <>
      <Header numPolls={polls.length} />
      <div className="wrapper">
        <div className="card frame">
          <PollList
            polls={polls}
            onDelete={handleDelete}
            onUpdate={handleUpdate}
          />
          {/* <SubmitForm onFormSubmit={this.handleSubmit} /> */}
        </div>
      </div>
    </>
  );
}

export default Home;
