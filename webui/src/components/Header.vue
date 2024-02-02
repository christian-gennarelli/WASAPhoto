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
                loading: true,
                showUpload: false,
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


        <div class="header-left" type="button" @click="this.$router.push('/home')"> 
            <span class="title"> WASAPhoto </span>
        </div>


        <div class="header-center">
            <input style="border-radius: 10px;" v-model="searchedUsername" @keyup.enter="this.$router.push({name: 'profile', params: {username: searchedUsername}}); searchedUsername=''" type="textbox" placeholder="Cerca qui...">
        </div>


        <div class="header-right"> 
            <router-link :to="{ name: 'profile', params: {username: user.Username} }"> <img :src="this.getImgUrl(user.ProfilePic)">{{user.Username}}</router-link> 
            <span type="button" @click="logout"> Logout </span>
        </div>

        <div v-if="showUpload" class="upload-overlay">
            <div class="upload-popup">

            </div>
        </div>

    </div>
</template>



<style scoped>

.header-grid-container {
    border: 2px solid black;
    border-radius: 10px;
    margin: 15px;
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
    /* background: radial-gradient(circle at 10% 20%, rgb(238, 56, 56) 0%, rgba(206, 21, 0, 0.92) 90.1%); */
    background: radial-gradient(circle at 10% 20%, rgb(255, 200, 124) 0%, rgb(252, 251, 121) 90%);
}

.header-left{
    margin: auto 10px;
}

.header-left .title {  
    font-size: 50px;
    font-weight: bold;
}

.header-center {
    margin: auto auto;
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
    font-size: 24px
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
    margin-left: 5px;
}

</style>