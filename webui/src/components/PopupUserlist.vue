<script>
    import { getImgUrl } from '../functions/getImgUrl'
    export default {
        props: ['list', 'category', 'show', 'headertxt', 'username'],
        emits: ['change-show', 'ban-user'],
        methods: {
            getImgUrl
        },
    }
</script>

<template>
    <div> 
        <span @click="this.$emit('change-show')" class="title"> {{ category }}</span>: {{ list.length }} 
        <div class="profile-overlay" v-if="show">
            <span class="profile-popup" >
                <span> {{ headertxt }} {{ username }}: </span>
                <img @click="this.$emit('change-show')" class="exit" src="@/assets/close.png">
                <div class="user" v-for="(user, key) in list" :key="key" >
                    <img class="profile-img" :src="this.getImgUrl(user.ProfilePic)" style="width: 30px; height: 30px;">
                    <router-link @click="this.$emit('change-show')" :to="{ name: 'profile', params: {username: user.Username }}">
                        {{ user.Username }}
                    </router-link>
                    <button v-if="category=='Banned'" class="btn-unban" @click="this.$emit('unban-user', user.Username)"> Unban </button>
                </div>
            </span>
        </div>
    </div>
</template>

<style scoped>

div {
    font-size: 25px;
    display: flex;
    align-items: center;
    justify-content: center;
}

div .title {
    font-weight: bold;
}

.profile-overlay .profile-popup a {
    text-decoration: none;
    color: black;
}

.profile-overlay .profile-popup a:hover {
    font-weight: bold;
}

.profile-overlay {

position: fixed;
top: 0;
bottom: 0;
left: 0;
right: 0;
background: black;
overflow:auto;
z-index: 1;

}

.profile-popup {

margin: 70px auto;
padding: 20px;
background: #fff;
border-radius: 5px;
width: 30%;
position: relative;
background: radial-gradient(circle at 10% 20%, rgb(255, 200, 124) 0%, rgb(252, 251, 121) 90%);

}

.profile-popup .user {
    display: block;
    margin-left: 0;
}

.profile-popup .profile-img {
    border-radius: 25px;
    margin-bottom: 5px;
    margin-right: 5px;
}

.profile-popup .btn-unban {
    border: none;
    background-color: transparent;
    font-size: 22px;
    height: 10px;
}

.profile-popup .btn-unban:hover {
    font-style: italic;
}

</style>