<script>
    import { getImgUrl } from '../functions/getImgUrl'
    import { logout } from '../functions/logout'
    export default {
        data() {
            return {
                user: {
                    username: '',
                    profilePic: '',
                },
                loading: true
            }
        },
        methods: {
            getImgUrl,
            logout
        },
        created(){
            this.user.username = localStorage.getItem('username')
            this.user.profilePic = localStorage.getItem('profilePic')
            this.loading = false
        }
    }

</script>

<template>
    <div class="custom">
        <div @click="this.$router.push('/home')"  style="font-size: 200%; font-weight: bold;"> WASAPhoto </div>
        <div style="display: inline; margin: auto">
            <input type="textbox" placeholder="Cerca qui... "> 
        </div>
        <div style="text-align: right"> 
            <img :src="this.getImgUrl(user.profilePic)" style="width: 32px; height: 32px"> 
            <router-link :to="{ name: 'profile', params: {username: user.username} }"> 
                {{ user.username }} 
            </router-link> 
            <button style="margin-left: 10px" @click="logout"> Logout </button>
        </div>
    </div>
</template>

<style scoped>

div.custom {
    display: grid;
    grid-template-columns: 34% 33% 33%;
    /* border: 2px solid red; */
}

</style>