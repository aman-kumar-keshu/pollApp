CREATE TABLE IF NOT EXISTS polls (
  id SERIAL NOT NULL PRIMARY KEY,
				name VARCHAR(255) UNIQUE NOT NULL,
				topic VARCHAR(255),
				src VARCHAR NOT NULL,
				upvotes INT NOT NULL,
				downvotes INT NOT NULL
);

INSERT INTO polls (name, topic, src, upvotes, downvotes) VALUES(
  'Angular','Awesome Angular', 'https://cdn.colorlib.com/wp/wp-content/uploads/sites/2/angular-logo.png', 1, 0
);

INSERT INTO polls(name, topic, src, upvotes, downvotes) VALUES(
  'Vue', 'Voguish Vue','https://upload.wikimedia.org/wikipedia/commons/thumb/5/53/Vue.js_Logo.svg/400px-Vue.js_Logo.svg.png', 1, 0
);

INSERT OR IGNORE INTO polls(name, topic, src, upvotes, downvotes) VALUES(
  'React','Remarkable React','https://upload.wikimedia.org/wikipedia/commons/thumb/a/a7/React-icon.svg/1200px-React-icon.svg.png', 1, 0);
	
INSERT OR IGNORE INTO polls(name, topic, src, upvotes, downvotes) VALUES(
  'Ember','Excellent Ember','https://cdn-images-1.medium.com/max/741/1*9oD6P0dEfPYp3Vkk2UTzCg.png', 1, 0);
	
INSERT OR IGNORE INTO polls(name, topic, src, upvotes, downvotes) VALUES(
  'Knockout','Knightly Knockout','https://images.g2crowd.com/uploads/product/image/social_landscape/social_landscape_1489710848/knockout-js.png', 1, 0);
	
