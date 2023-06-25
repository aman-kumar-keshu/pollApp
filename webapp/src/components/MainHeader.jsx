import {MdPostAdd, MdMessage} from 'react-icons/md';

function MainHeader({onCreatePoll}){
    return (
        <header>
            <h1 >
                <MdMessage/>  
                    Poll Company
            </h1>

            <p>
                <button onClick = {onCreatePoll}>
                    <MdPostAdd size = {18}/> 
                </button>
            </p>
        </header>
    );
}

export default MainHeader;