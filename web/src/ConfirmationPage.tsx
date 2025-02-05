import { useNavigate, useParams } from "react-router";
import { API_URL } from "./App";

function ConfirmationPage() {
    const {token = ''} = useParams();
    const navigate = useNavigate();

    const handleConfirm = async () => {
        const response = await fetch(`${API_URL}/auth/activate/${token}`, {
            method: 'PUT'
        });

        if (response.ok) {
            navigate('/');
        } else {
            // handle err 
            alert('Account not activated');
        }
    }

    return (
        <div>
            <h1>Confirmation</h1>
            <button onClick={handleConfirm}>Click to confirm</button>
        </div>
    )
}

export default ConfirmationPage;