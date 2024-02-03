<script>

export default {
    data() {
        return {
            username: '',
        }
    }, 
    methods: {
        async doLogin() {
            await this.$axios.post(
                "/session",
                this.username, // Username passed in plain text
                { // Request headers
                    headers: {
                        'Content-Type': 'text/plain'
                    },
                    responseType: 'text',
                }
            ).then((res) => {
                console.log(res.data)
                localStorage.setItem('Username', this.username) // Store the input username in local storage
                localStorage.setItem('ID', res.data.ID) // Store token from response in local storage
                localStorage.setItem('Birthdate', res.data.Birthdate)
                localStorage.setItem('Name', res.data.Name)
                localStorage.setItem('ProfilePic', res.data.ProfilePic)
                this.$router.push('/home') // Route to home page
            }).catch((e) => {
                let error = e.response.data
                alert(error.ErrorCode + " " + error.Description)
            })

        }
    }
}

</script>

<template>
    <div class="login-container">
        <span class="login-title"> WASAPhoto. </span>
        <div class="login-form">
            <h2> Get started! </h2>
            <input type="text" class="login-input" v-model="username" @keyup.enter="doLogin" placeholder="Enter username">
            <button class="login-btn" @click.prevent="doLogin"> Go! </button>
        </div>
    </div>
</template>

<style scoped>

.login-container {
    display: grid;
    grid-template-columns: 3fr 1fr;
    grid-template-rows: 100vh;
    align-items: center;
    justify-content: center;
}

.login-title {
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: bold;
    font-size: 10em;
}

.login-form {
    width: 100%;
}

.login-input {
    outline: none;
    border: 1px solid black;
    background: transparent;
}

::placeholder {
    color: black
}

.login-btn {
    margin-left: 5px;
    background: transparent;
    border: 1px solid black;
    outline: none;
}

.login-btn:hover {
    border: 3px solid black;
    font-weight: bold;
}

</style>