import {createRouter, createWebHashHistory} from 'vue-router'
import Home from '../views/HomeView.vue'
import Login from '../views/LoginView.vue'
import NotFound from '../views/NotFoundView.vue'
import Profile from '../views/ProfileView.vue'


const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{
			path: '/',
			name: 'login',
			component: Login
		},
		{
			path: '/home', 
			name: 'home',
			component: Home,
		},
		{
			path: '/users/:username/profile', 
			name: 'profile',
			component: Profile,
		},
		// Redirection for non-existent paths
		{
			path: '/404',
			name: 'NotFound',
			component: NotFound
		},
		{
			path: '/:catchAll(.*)',
			redirect: '/404'
		}
	]
})

// Redirect if:
router.beforeEach(async (to, from) => {

	// the non-authenticated client asks for a page different from the login one
	let token = localStorage.getItem('ID');
	if (token === null && to.name !== 'login') {
		return { name: 'login' };
	}

	// the authenticated client asks for the login page
	if (token !== null && to.name === 'login') {
		return { name: 'home' };
	}

});

export default router
