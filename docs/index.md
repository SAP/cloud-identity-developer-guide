---
layout: home

hero:
  name: "AMS Client Libraries"
  text: "Documentation"
  tagline: "SAP Cloud Identity Services Authorization"
  actions:
    - theme: brand
      text: Getting Started
      link: /concepts/GettingStarted
    - theme: alt
      text: Samples
      link: /Samples

features:
  - icon: 
      src: /declare.svg
      alt: Declare
    title: Declare
    details: Use Data Control Language (DCL) to declare the actions, resources and their attributes on which the authorization model of your business application should be based on. The result is a set of base policies for your business application.
  - icon: 
      src: /enforce.svg
      alt: Enforce
    title: Enforce
    details: Use the AMS Client Libraries to perform authorization checks, so that the security relevant operations and resources of your business application are protected from unauthorized access.
    class: feature-highlight
  - icon: 
      src: /manage.svg
      alt: Manage
    title: Manage
    details: Create administration policies from base policies and assign them to the users of your business application, so that they can only perform those actions on the resources for which they have been authorized for.
---