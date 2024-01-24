<script>

export default {
    data() {
        return {
            username: '',
        }
    }, 
    methods: {
        async doLogin() {
            this.$axios.post(
                "/session",
                this.username, // Username passed in plain text
                { // Request headers
                    headers: {
                        'Content-Type': 'text/plain'
                    },
                    responseType: 'text',
                }
            ).then((res) => {
                localStorage.setItem('username', this.username) // Store the input username in local storage
                localStorage.setItem('token', res.data.ID) // Store token from response in local storage
                localStorage.setItem('birthdate', res.data.Birthdate)
                localStorage.setItem('name', res.data.Name)
                localStorage.setItem('profilePic', res.data.ProfilePic)
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
    <div>
        <h1> WASAPhoto </h1>
        <p class="center">
            <label> Username <input type="text" v-model="username" @keyup.enter="doLogin" placeholder="Click here to insert"> </label>
            <button class="right" @click.prevent="doLogin"> Login </button>
        </p>
    </div>
</template>