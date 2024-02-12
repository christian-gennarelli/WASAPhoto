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
    <div style="font-size: 25px;"> 
        <span type="button" @click="this.$emit('change-show')" style="font-weight: bold;"> {{ category }}:</span> {{ list.length }}
        <div class="profile-overlay" v-if="show">
            <span class="profile-popup" style="font-size: 25px;" >
                <span> {{ category }}: </span>
                <img type="button" @click="this.$emit('change-show')" class="exit" src="@/assets/buttons/close.png">
                <div class="user" v-for="(user, key) in list" :key="key" >
                    <img class="profile-img" :src="this.getImgUrl(user.ProfilePic)" style="width: 30px; height: 30px;">
                    <span v-if="this.category == 'Banned' "> {{ user.Username }} </span>
                    <router-link v-else @click="this.$emit('change-show')" :to="{ name: 'profile', params: {username: user.Username }}">
                        {{ user.Username }}
                    </router-link>
                    <button v-if="category=='Banned'" class="banned-btn" @click="this.$emit('unban-user', user.Username)"> Unban </button>
                </div>
            </span>
        </div>
    </div>
</template>

<style scoped>

.banned-btn {
    background-color: transparent;
    border: none;
    font-size: 18px;
    margin-right: 2px;
}

.banned-btn:hover{
    color: green;
    font-style: italic;
}

.profile-overlay .profile-popup a {
    text-decoration: none;
    color: black;
}

.profile-overlay .profile-popup a:hover {
    font-weight: bold;
}

.profile-overlay {
    display: flex;
    align-items: center;
    justify-content: center;
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    background: black;
    overflow: auto;
    z-index: 1;
}

.profile-popup {

    padding: 20px;
    background: #fff;
    border-radius: 5px;
    width: auto;
    overflow: scroll;
    min-width: 400px;
    position: relative;
    background: radial-gradient(circle at 10% 20%, rgb(255, 200, 124) 0%, rgb(252, 251, 121) 90%);
    z-index: 1;

}

.profile-popup::-webkit-scrollbar {
  display: none;
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