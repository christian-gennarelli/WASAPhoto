<script>
    import Post from '../components/Post.vue'
    import PopupUserlist from '../components/PopupUserlist.vue'
    import { getImgUrl } from '../functions/getImgUrl'
    import { logout } from '../functions/logout'

    export default {
        data() {
            return {
                visitedProfile: {
                    user: {
                        Username: '',
                        Name: '',
                        Birthdate: '',
                        ProfilePic: '',
                    },
                    posts: [],
                    followers: [],
                    followings: [],
                    banned: [],
                },
                user: {
                    ID: localStorage.getItem('ID'),
                    Username: localStorage.getItem('Username'),
                    followings: [],
                },
                newUsername: '',
                loading: true,
                modifying: false,
                followed: false,
                showFollowers: false,
                showFollowings: false,
                showBanned: false,
            }
        },
        methods: {
            updateUsername(){
                this.$axios.put(
                    "/users/" + this.user.Username + "/profile/",
                    this.newUsername,
                    {
                        headers: {
                            "Authorization": this.user.ID,
                            "Content-Type": 'text/plain'
                        }
                    }
                ).then(()=>{
                    alert("Username updated! New username: " + this.newUsername)
                    this.newUsername = ''
                    this.loading = true
                    this.modifying = false
                    this.logout()
                }).catch((e)=>{
                    alert(e.response.data.ErrorCode + " " + e.response.data.Description)
                })
            },
            followUser(){
                this.$axios.put(
                    '/users/' + this.user.Username + '/followings/',
                    null,
                    {
                        params: {
                            followed_username: this.visitedProfile.user.Username,
                        },
                        headers: {
                            'Authorization': this.user.ID
                        }
                    }
                ).catch(()=>{
                    alert(e.response.data.ErrorCode + " " + e.response.data.Description)
                })
                this.visitedProfile.followers.push(this.user)
                this.followed = !this.followed 
            },
            unfollowUser(){
                this.$axios.delete(
                    '/users/' + this.user.Username + '/followings/' + this.visitedProfile.user.Username,
                    {
                        headers: {
                            "Authorization": this.user.ID,
                        }
                    }
                ).catch(()=>{
                    alert(e.response.data.ErrorCode + " " + e.response.data.Description)
                })
                this.visitedProfile.followers = this.visitedProfile.followers.filter(u => u.Username !== this.user.Username)
                this.followed = !this.followed 
            },
            banUser(){
                this.$axios.put(
                    '/users/' + this.user.Username + '/banned/',
                    null,
                    {
                        params: {
                            banned_username: this.visitedProfile.user.Username
                        },
                        headers: {
                            'Authorization': this.user.ID
                        }
                    }
                ).catch((e)=>{
                    alert(e.response.data.ErrorCode + " " + e.response.data.Description)
                })
                alert("You banned this user successfully! Press OK to go back on the home page...")
                this.$router.push('/home')
            },
            unbanUser(bannedUsername){
                this.$axios.delete(
                    '/users/' + this.user.Username + '/banned/' + bannedUsername,
                    {
                        headers: {
                            'Authorization': this.user.ID
                        }
                    }
                ).catch((e)=>{
                    alert(e.response.data.ErrorCode + " " + e.response.data.Description)
                })
                this.visitedProfile.banned = this.visitedProfile.banned.filter(u => u.Username !== bannedUsername)
            },
            async getVisitedProfile(){
                // GetUserProfile (visited user)
                await this.$axios.get(
                    "/users/" + this.$route.params.username + "/profile/",
                    {
                        headers: {
                            'Authorization': this.user.ID   
                        },
                    }
                ).then((res) => {
                    this.visitedProfile.user.Username = res.data.User.Username
                    this.visitedProfile.user.Name = res.data.User.Name
                    this.visitedProfile.user.Birthdate = res.data.User.Birthdate
                    this.visitedProfile.user.ProfilePic = res.data.User.ProfilePic
                    this.visitedProfile.followings = res.data.Followings ? res.data.Followings : []
                    this.visitedProfile.followers = res.data.Followers ? res.data.Followers : []
                    this.visitedProfile.banned = res.data.Banned ? res.data.Banned : []
                    this.visitedProfile.posts = res.data.Posts ? res.data.Posts : []
                    console.log(this.visitedProfile.posts)
                }).catch((e) => {
                    alert(e.response.data.ErrorCode + " " + e.response.data.Description)
                    this.$router.push('/')
                })
                // GetFollowingsList (authenticated user)
                await this.$axios.get(
                    "/users/" + this.user.Username + "/followings/",
                    {
                        headers: {  
                            'Authorization': localStorage.getItem('ID')
                        }
                    }
                ).then((res)=>{
                    if (this.visitedProfile.user.Username !== this.user.Username) {
                        for (let i = 0; i < res.data.length; i++){
                            if (res.data[i]["Username"] === this.visitedProfile.user.Username){
                                this.followed = true
                                break
                            }
                        }
                    }
                }).catch((e)=>{
                    alert(e.response.data.ErrorCode + " " + e.response.data.Description)
                    this.loading = true
                    this.$router.push('/')
                })
                // Update vars
                this.loading = false
            },
            getImgUrl,
            logout
        },
        watch:{
            $route (to, from){
                if (to.name == "profile") {
                    this.loading = true
                    this.followed = false
                    this.getVisitedProfile()
                }
            }
        },
        async created() { 
            this.getVisitedProfile()
        },
        components:{
            Post: Post, PopupUserlist
        }
    }
</script>

<template>
    <Header></Header>
    <div v-if="!loading" class="profile-container"> 
        <div class="left">
            <img class="profile-img" :src="this.getImgUrl(this.visitedProfile.user.ProfilePic)">
            <span class="profile-username" v-if="!this.modifying"> {{ visitedProfile.user.Username }} </span>
            <span class="modified-profile-username" v-if="this.modifying">
                <input type="textbox" v-model="newUsername" @keyup.enter="updateUsername" placeholder="Enter new username">
                <!-- <span> Remember: username can contain only lower/upper case letters, underscores (_) and dashes (-). </span> -->
            </span>
            <img class="pencil" src="@/assets/pencil.png" v-if="this.visitedProfile.user.Username === this.user.Username" @click="this.modifying=!this.modifying">
            <button v-if="followed" class="btn-red" button @click="unfollowUser"> Unfollow </button>
            <button v-else-if="this.visitedProfile.user.Username !== this.user.Username" class="btn-green" @click="followUser"> Follow </button>
            <button v-if="this.visitedProfile.user.Username !== this.user.Username" class="btn-red" @click="banUser"> Ban </button>
        </div>

        <div class="right"> 
            <PopupUserlist
                category="Followers"
                headertxt="Users following "
                :username="this.visitedProfile.user.Username"
                :show="this.showFollowers"
                :list="this.visitedProfile.followers"
                @change-show="this.showFollowers=!this.showFollowers; this.showFollowings=false; this.showBanned=false"
            >
            </PopupUserlist>
            <PopupUserlist
                category="Followings"
                headertxt="Users followed by "
                :username="this.visitedProfile.user.Username"
                :show="this.showFollowings"
                :list="this.visitedProfile.followings"
                @change-show="this.showFollowings=!this.showFollowings; this.showFollowers=false; this.showBanned=false"
            >
            </PopupUserlist>
            <PopupUserlist
                category="Banned"
                headertxt="Users banned by "
                :username="this.visitedProfile.user.Username"
                :show="this.showBanned"
                :list="this.visitedProfile.banned"
                @change-show="this.showBanned=!this.showBanned; this.showFollowers=false; this.showFollowings=false"
                @unban-user="(bannedUsername) => unbanUser(bannedUsername)"
            >
            </PopupUserlist>
        </div>
    </div>
    
    <div class="home-container" v-if="!showFollowers && !showFollowings && !showBanned">
        <Post  
            v-if="visitedProfile.posts"
            v-for="(post, key) in visitedProfile.posts" 
            :post="post"
            :user="user"
            @add-like="this.visitedProfile.posts[key].Likes ? this.visitedProfile.posts[key].Likes.push(this.user) : this.visitedProfile.posts[key].Likes = [user]"
            @remove-like="this.visitedProfile.posts[key].Likes = this.visitedProfile.posts[key].Likes.filter(u => u.Username !== this.user.Username)"
            @remove-post="this.visitedProfile.posts.splice(key, 1)"
        ></Post>
    </div>
</template>

<style>

.profile-container {
    border: 2px solid black;
    margin: 15px;
    border-radius: 15px;
    display: grid;
    grid-template-columns: 50% 50%;
    height: 120px;
}

.left {
    display: flex;
    justify-content: center;
    align-items: center;
}

.left button {
    border: 2px solid black;
    background-color: transparent;
    font-size: 30px;
    margin-top: 5px;
    margin-right: 2px;
}

.left .btn-red:hover {
    background: linear-gradient(108.4deg, rgb(253, 44, 56) 3.3%, rgb(176, 2, 12) 98.4%);
}

.left .btn-green:hover {
    background: linear-gradient(to right, rgb(182, 244, 146), rgb(51, 139, 147));
}

.left .profile-username {
    font-size: 70px;
}

.left .profile-img {
    width: 70px;
    height: 70px;
    border-radius: 50px;
    margin-right: 15px;
}

.left input {
    height: 60px;
}

.left .pencil {
    width: 50px;
    height: 50px;
    margin-left: 15px;
}

.left .pencil:hover {
    cursor: pointer;
}

</style>