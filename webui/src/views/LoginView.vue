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
        <div class="login-title"> WASAPhoto. </div>
        <div class="login-form">
            <h2> Get started!</h2>
            <input type="login-input" v-model="username" @keyup.enter="doLogin" placeholder="Username">
            <button class="login-btn" @click.prevent="doLogin"> Go! </button>
        </div>
    </div>
</template>

<style scoped>

.login-container {
    display: grid;
    position: absolute;
    width: 1100px;
    grid-template-columns: 50fr 50fr;
    margin: 0;
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
}

.login-title {
    font-size: 150px;
}

.login-form {
    margin: auto;
}

.login-input{
    border: none;
    border-bottom: 1px solid black;
    outline: none;
}

.login-btn{
    margin-left: 5px;
}

</style>