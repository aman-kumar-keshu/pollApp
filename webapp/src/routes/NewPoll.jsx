import { useState } from "react";
import { createPoll } from "../services/pollService";
import { useNavigate } from "react-router-dom";

function NewPoll(props) {
  const [name, setName] = useState("");
  const [src, setSrc] = useState("");
  const [topic, setTopic] = useState("");
  const [options, setOptions] = useState([]);

  const [newOption, setNewOption] = useState("");

  const addOption = (newOption) => {
    setOptions((existingOptions) => [...existingOptions, newOption]);
    setNewOption("");
  };

  const navigator = useNavigate();

  const handleNewPostSubmit = async (e) => {
    e.preventDefault();
    console.log("add new poll callback", name, src, topic);
    const poll = { name, src, topic };
    await createPoll(poll);
    navigator("/");
  };

  return (
    <form onSubmit={handleNewPostSubmit}>
      <label>Poll Name</label>
      <input
        type="text"
        value={name}
        onChange={(e) => setName(e.target.value)}
      />
      <label> Poll Topic</label>

      <input
        type="text"
        value={topic}
        onChange={(e) => setTopic(e.target.value)}
      />
      <label> Image Url</label>

      <input type="text" value={src} onChange={(e) => setSrc(e.target.value)} />

      <ul>
        {options.map((option) => (
          <li>{option}</li>
        ))}{" "}
      </ul>
      <div>
        <input
          type="text"
          value={newOption}
          onChange={(e) => setNewOption(e.target.value)}
        />
        <button type="button" onClick={() => addOption(newOption)}>
          +
        </button>
      </div>
      <button type="submit"> Submit </button>
    </form>
  );
}

export default NewPoll;
// {
/* <ul id="mainMenu" role="menubar">
  <li class="menu-item" role="presentation">
    <a role="menuitem" tabindex="0" data-turbo="false" href="/en/users/sign_up">
      <div class="puIcon puIcon-user" aria-hidden="true"></div>
      <div class="text" role="presentation">
        Register
      </div>
    </a>
  </li>

  <li class="menu-item" role="presentation">
    <a role="menuitem" tabindex="-1" href="/en/users/sign_in">
      <div class="puIcon puIcon-login" aria-hidden="true"></div>
      <div class="text" role="presentation">
        Login
      </div>
    </a>
  </li>

  <li class="menu-item" role="presentation">
    <a role="menuitem" tabindex="-1" href="/en/accounts">
      <div class="puIcon puIcon-price" aria-hidden="true"></div>
      <div class="text" role="presentation">
        Pricing
      </div>
    </a>
  </li>
  <li class="menu-item" role="presentation">
    <a
      role="menuitem"
      aria-haspopup="true"
      aria-expanded="false"
      tabindex="-1"
      href="/en/support"
    >
      <div class="puIcon puIcon-support" aria-hidden="true"></div>
      <div class="text" role="presentation">
        Support
      </div>
    </a>
    <ul aria-label="Support" role="menu">
      <li class="subMenuItem-onlyExpanded" role="presentation">
        <a role="menuitem" href="/en/support">
          <div class="puIcon puIcon-support" aria-hidden="true"></div>
          <div class="text" role="presentation">
            Support
          </div>
        </a>
      </li>
      <li role="presentation">
        <a role="menuitem" href="/en/tutorials">
          <div class="puIcon puIcon-tutorials" aria-hidden="true"></div>
          <div class="text" role="presentation">
            Tutorials
          </div>
        </a>
      </li>
      <li role="presentation">
        <a role="menuitem" href="/en/posts">
          <div class="puIcon puIcon-blog" aria-hidden="true"></div>
          <div class="text" role="presentation">
            Blog
          </div>
        </a>
      </li>
      <li role="presentation">
        <a role="menuitem" href="/en/help">
          <div class="puIcon puIcon-info" aria-hidden="true"></div>
          <div class="text" role="presentation">
            Guide
          </div>
        </a>
      </li>
      <li role="presentation">
        <a role="menuitem" data-turbo="false" href="/forum?locale=en">
          <div class="puIcon puIcon-comments" aria-hidden="true"></div>
          <div class="text">Support Forum</div>
        </a>
      </li>
    </ul>
  </li>
  <li class="buildHelperButton" role="presentation">
    <a
      rel="nofollow"
      class="btn build-helper"
      role="menuitem"
      data-remote="true"
      href="/en/polls/build_helper?build_type=poll_type"
    >
      Create PollUnit
    </a>
  </li>
</ul>; */
