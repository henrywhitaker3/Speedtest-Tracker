# Contributing

✨ Thank you for contributing ✨

Please feel free to contribute by submitting PR's with your own helper, improvement to code snippets, explanations, etc.


## Submitting an issue

Found a problem? Have an enhancement? 

First of all see if your issue or idea has already been [reported](https://github.com/shystruk/famulus/issues).

If do not, open a [new one](https://github.com/shystruk/famulus/issues/new).


## Add your own helper

- Please take for example helper **isValuesUnique.js**
- Create *[helper_name].js* file in the root folder
- In docs folder create *[helper_name].md* file and fill in data keeping the structure
- In tests folder create *[helper_name].test.js* file and cover the helper by unit tests and much as possible
- In famulus.js file import and export already created helper
- In documentation section at README.md file add the helper to a specific category. If the category does not exist add the new one
- Before submitting a PR please check *Submitting a pull request* section below


## Submitting a pull request

- Fork this repository
- Clone fork `git clone ...`
- Navigate to the cloned directory
- Install all dependencies `npm install`
- Crate a new branch for the feature `git checkout -b new-feature`
- Make changes
- Run tests `npm run test`
- Commit changes `git commit -am 'What is feature about? :)'`
- Push to the branch `git push origin new-feature`
- Submit a PR
