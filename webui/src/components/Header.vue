<script>
    import { getImgUrl } from '../functions/getImgUrl'
    import { logout } from '../functions/logout'
    export default {
        data() {
            return {
                user: {
                    Username: '',
                    ProfilePic: '',
                },
                searchedUsername: '',
                loading: true
            }
        },
        methods: {
            getImgUrl,
            logout
        },
        created(){
            this.user.Username = localStorage.getItem('Username')
            this.user.ProfilePic = localStorage.getItem('ProfilePic')
            this.loading = false
        }
    }

</script>

<template>
    <div class="header-grid-container">


        <div class="header-left" @click="this.$router.push('/home')"> 
            <span> WASAPhoto </span>
        </div>


        <div class="header-center">
            <input v-model="searchedUsername" @keyup.enter="this.$router.push({name: 'profile', params: {username: searchedUsername}}); searchedUsername=''" type="textbox" placeholder="Cerca qui...">
        </div>


        <div class="header-right"> 
            <router-link :to="{ name: 'profile', params: {username: user.Username} }"> <img :src="this.getImgUrl(user.ProfilePic)">{{user.Username}}</router-link> 
            <span @click="logout"> Logout </span>
        </div>

    </div>
</template>



<style scoped>

.header-grid-container {
    border: 2px solid black;
    border-radius: 10px;
    margin: 15px;
    display: grid;
    grid-template-columns: 34fr 33fr 33fr;
    background: radial-gradient(circle at 10% 20%, rgb(238, 56, 56) 0%, rgba(206, 21, 0, 0.92) 90.1%);
}

.header-left{
    margin: auto 10px;
}

.header-left span {  
    font-size: 50px;
    font-style: bold;
}

.header-left span:hover {
    text-decoration: underline;
    text-decoration-thickness: 4px;
}

.header-center {
    margin: auto auto;
}

.header-center input {
    border-radius: 10px;
}

.header-right{
    margin-left: auto;
    margin-right: 5%;
    margin-top: auto;
    margin-bottom: auto;
}

.header-right a {
    font-size: 20px;
    font-weight: bold;
    text-decoration: none;
    color: black
}

.header-right a:hover {
    text-decoration: underline;
}

.header-right img {
    margin-bottom: 3%;
    border-radius: 25px;
    width: 32px;
    height: 32px;
}

.header-right span {
    font-size: 20px;
    font-style: italic;
}

.header-right span:hover {
    cursor: pointer;
}


</style>