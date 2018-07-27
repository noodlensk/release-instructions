console.log(process.env.NODE_ENV)
export default {
    API_URL: process.env.NODE_ENV === 'development' ? 'http://localhost:8081' : '',
    ENDPOINT_MAP: {
        
    }
};
